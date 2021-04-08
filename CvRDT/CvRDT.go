package CvRDT

import (
	"encoding/json"
	"fmt"
	"github.com/er1c-zh/crdt-viz/network"
	"sync"
	"sync/atomic"
	"time"
)

type State interface {
	Init(opt InitOption)
	Query() interface{}
	Update(args UpdateArgs)
	Merge(state State)
	GetStatus() json.Marshaler
}

type InitOption struct {
	Id int
}

type UpdateArgs struct {
}

type Option struct {
	Count int
}

type StateFactory func() State

type Obj struct {
	Factory StateFactory
	Opt     Option
	Cloud   network.Cloud

	instMap      map[int]State
	interfaceMap map[string] /*addr*/ network.Interface
	msgWaitSend  sync.Map
	serial       int64
}

func NewObject(factory StateFactory) *Obj {
	return &Obj{
		Factory:      factory,
		Opt:          Option{Count: 10},
		instMap:      map[int]State{},
		interfaceMap: map[string]network.Interface{},
	}
}

const (
	ActionDoMerge = iota + 1
	ActionDoneMerge
)

type broadcastMsg struct {
	Action int   `json:"action"`
	Serial int64 `json:"serial"`
	From   int   `json:"from"`
	To     int   `json:"to"`
}

func (o *Obj) Init(cloud network.Cloud, opt Option) {
	getAddr := func(id int) string {
		return fmt.Sprintf("%d", id)
	}
	o.Opt = opt

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				o.msgWaitSend.Range(func(_, value interface{}) bool {
					v := value.(broadcastMsg)
					j, _ := json.Marshal(v)
					_ = o.interfaceMap[getAddr(v.From)].Write(getAddr(v.To), string(j))
					return true
				})
			}
		}
	}()

	o.Cloud = cloud
	for i := 0; i < opt.Count; i++ {
		o.instMap[i] = o.Factory()
		o.instMap[i].Init(InitOption{
			Id: i,
		})
		addr := getAddr(i)
		_interface, err := o.Cloud.NewClient(addr, network.DefaultOption())
		if err != nil {
			panic("create interface fail")
		}
		o.interfaceMap[addr] = _interface
		ch, err := _interface.GetRecv()
		if err != nil {
			panic("get ch fail")
		}
		// listen msg
		go func() {
			for {
				select {
				case msg := <-ch:
					m := &broadcastMsg{}
					err = json.Unmarshal([]byte(msg.Msg), m)
					if err != nil {
						// todo log
						continue
					}
					switch m.Action {
					case ActionDoMerge:
						o.instMap[m.To].Merge(o.instMap[m.From])
						go func() {
							// send ack
							ack := broadcastMsg{
								Action: ActionDoneMerge,
								Serial: m.Serial,
								From:   m.To,
								To:     m.From,
							}
							j, _ := json.Marshal(ack)
							_ = o.interfaceMap[getAddr(m.To)].Write(getAddr(m.From), string(j))
						}()
					case ActionDoneMerge:
						o.msgWaitSend.Delete(m.Serial)
					}
				}
			}
		}()
	}

	for _, i := range o.interfaceMap {
		_ = i.Enable()
	}
}

func (o *Obj) Query(Idx int) interface{} {
	return o.instMap[Idx].Query()
}

func (o *Obj) Update(Idx int, args UpdateArgs) {
	o.instMap[Idx].Update(args)
	// broadcast
	for i := 0; i < o.Opt.Count; i++ {
		if i == Idx {
			continue
		}
		serial := atomic.AddInt64(&o.serial, 1)
		msg := broadcastMsg{
			Action: ActionDoMerge,
			Serial: serial,
			From:   Idx,
			To:     i,
		}
		o.msgWaitSend.Store(serial, msg)
	}
}

func (o *Obj) Merge(Idx int, FromIdx int) {
	o.instMap[Idx].Merge(o.instMap[FromIdx])
}

func (o *Obj) State() interface{} {
	m := map[int]json.Marshaler{}
	for id, s := range o.instMap {
		m[id] = s.GetStatus()
	}
	return m
}
