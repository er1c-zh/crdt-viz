package network

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	ChaosActionSplitBrain = iota + 1
	ChaosActionLoss
)

type ChaosAction struct {
	Action int
	Args   string
}

type Filter = func(pkg MsgPkg) bool /*是否被过滤*/

type FilterFactory = func(action ChaosAction) (Filter, error)

var (
	FilterFactoryMap = map[int]FilterFactory{
		ChaosActionSplitBrain: func(action ChaosAction) (Filter, error) {
			addrList := strings.Split(action.Args, ",")
			m := map[string]interface{}{}
			for _, addr := range addrList {
				m[addr] = struct{}{}
			}
			return func(pkg MsgPkg) bool {
				ok1 := m[pkg.From]
				ok2 := m[pkg.To]
				if ok1 != ok2 {
					return true
				}
				return false
			}, nil
		},
		ChaosActionLoss: func(action ChaosAction) (Filter, error) {
			argList := strings.Split(action.Args, ",")
			if len(argList) != 3 {
				return nil, errors.New("arg should be addr1,addr2,loss_pct")
			}
			addr1 := argList[0]
			addr2 := argList[1]
			pctStr := argList[2]
			pct, err := strconv.ParseInt(pctStr, 10, 64)
			if err != nil {
				return nil, errors.New("loss_pct should be integer")
			}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			return func(pkg MsgPkg) bool {
				if (pkg.From == addr1 && pkg.To == addr2) ||
					(pkg.From == addr2 && pkg.To == addr1) {
					return r.Int63n(100) < pct
				}
				return false
			}, nil
		},
	}
)
