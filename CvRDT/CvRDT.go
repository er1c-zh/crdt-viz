package CvRDT

import "encoding/json"

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

	instMap map[int]State
}

func NewObject(factory StateFactory) *Obj {
	return &Obj{
		Factory: factory,
		Opt:     Option{Count: 10},
		instMap: map[int]State{},
	}
}

func (o *Obj) Init(opt Option) {
	for i := 0; i < opt.Count; i++ {
		o.instMap[i] = o.Factory()
		o.instMap[i].Init(InitOption{
			Id: i,
		})
	}
}

func (o *Obj) Query(Idx int) interface{} {
	return o.instMap[Idx].Query()
}
func (o *Obj) Update(Idx int, args UpdateArgs) {
	o.instMap[Idx].Update(args)
	// todo broadcast
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
