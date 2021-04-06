package network

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 一个模拟分布式网络的工具

type Logger interface {
	Logf(string, ...interface{})
}

func NewCloud(opt ...Option) Cloud {
	return newCloud(opt...)
}

type Interface interface {
	GetRecv() (<-chan MsgPkg, error)
	Write(addr string, msg string) error
	Enable() error
}

type Cloud interface {
	NewClient(addr string, option ...Option) (Interface, error)
}

type Option struct {
	BufLenRecv                      int
	BufLenWrite                     int
	LatencyInMillisecond            int64
	LatencyRandomDeltaInMillisecond int64
	PkgLossPct                      int
}

func DefaultOption() Option {
	return Option{
		BufLenRecv:                      10,
		BufLenWrite:                     10,
		LatencyInMillisecond:            10,
		LatencyRandomDeltaInMillisecond: 3,
		PkgLossPct:                      1,
	}
}

////////////////////////////////////////////
// private
////////////////////////////////////////////

type logger struct{}

func (l logger) Logf(s string, i ...interface{}) {
	fmt.Printf("[mock_network]"+s+"\n", i...)
}

type express interface {
	Send(msg MsgPkg) error
	Enable(*client) error
}

func newClient(inst express, addr string, option Option) *client {
	return &client{
		addr:    addr,
		recv:    make(chan MsgPkg, option.BufLenRecv),
		write:   make(chan MsgPkg, option.BufLenWrite),
		express: inst,
	}
}

type MsgPkg struct {
	From string
	To   string
	Msg  string

	SendTime             int64
	expectedDeliveryTime int64
}

type client struct {
	addr  string
	recv  chan MsgPkg
	write chan MsgPkg
	express
}

func (c *client) GetRecv() (<-chan MsgPkg, error) {
	return c.recv, nil
}

func (c *client) Write(addr string, msg string) error {
	return c.Send(MsgPkg{
		From:     c.addr,
		To:       addr,
		Msg:      msg,
		SendTime: time.Now().UnixNano(),
	})
}

func (c *client) Enable() error {
	return c.express.Enable(c)
}

type cloud struct {
	net    sync.Map
	bus    chan MsgPkg
	option Option
	Logger
}

func newCloud(opt ...Option) *cloud {
	inst := &cloud{
		option: DefaultOption(),
		Logger: logger{},
		bus:    make(chan MsgPkg),
	}
	rand.Seed(time.Now().UnixNano())
	if len(opt) > 0 {
		inst.option = opt[0]
	}
	go inst.loop()
	return inst
}

func (c *cloud) loop() {
	queue := make([]MsgPkg, 0)
	ticker := time.NewTicker(1 * time.Millisecond)
	delivery := func(msg *MsgPkg) {
		if time.Now().UnixNano()-msg.expectedDeliveryTime < 0 {
			// skip
			queue = append(queue, *msg)
			return
		}
		to, ok := c.net.Load(msg.To)
		if !ok {
			// skip
			return
		}
		select {
		case to.(*client).recv <- *msg:
			// done
		default:
			// skip
			queue = append(queue, *msg)
		}
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				// try one old
				if len(queue) == 0 {
					continue
				}
				delivery(&(queue[0]))
				queue = queue[1:]
			}
		}
	}()

	for {
		select {
		case msg := <-c.bus:
			delivery(&msg)
		}
	}
}

func (c *cloud) NewClient(addr string, option ...Option) (Interface, error) {
	opt := c.option
	if len(option) > 0 {
		opt = option[0]
	}
	tmp := newClient(c, addr, opt)
	return tmp, nil
}

func (c *cloud) Send(msg MsgPkg) error {
	// 寻址
	if _, ok := c.net.Load(msg.To); !ok {
		c.Logf("Send not found")
		return errors.New("target not found")
	}
	msg.expectedDeliveryTime = msg.SendTime +
		(c.option.LatencyInMillisecond+
			(rand.Int63n(2*c.option.LatencyRandomDeltaInMillisecond)-c.option.LatencyRandomDeltaInMillisecond))*
			int64(time.Millisecond)
	if c.option.PkgLossPct > 0 && rand.Intn(100) < c.option.PkgLossPct {
		c.Logf("mock loss!")
		return nil
	}
	c.bus <- msg
	return nil
}

func (c *cloud) Enable(client *client) error {
	addr := client.addr
	_, ok := c.net.LoadOrStore(addr, client)
	if ok {
		return errors.New("addr already register")
	}
	return nil
}
