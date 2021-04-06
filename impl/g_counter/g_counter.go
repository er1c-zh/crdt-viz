package g_counter

import (
	"encoding/json"
	"github.com/er1c-zh/crdt-viz/CvRDT"
)

var Factory CvRDT.StateFactory = func() CvRDT.State {
	return &State{}
}

type State struct {
	Id          int
	CounterList []int
}

func (s *State) Init(opt CvRDT.InitOption) {
	s.Id = opt.Id
	s.CounterList = make([]int, 0)
}

func (s *State) Query() interface{} {
	cnt := 0
	for _, c := range s.CounterList {
		cnt += c
	}
	return cnt
}

func (s *State) Update(args CvRDT.UpdateArgs) {
	for len(s.CounterList) <= s.Id {
		s.CounterList = append(s.CounterList, 0)
	}
	s.CounterList[s.Id]++
}

func (s *State) Merge(state CvRDT.State) {
	_s := state.(*State)
	for i, c := range _s.CounterList {
		s.CounterList[i] += c
	}
}

type jsonWrapper struct {
	State
}

func (j jsonWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

func (s *State) GetStatus() json.Marshaler {
	return jsonWrapper{*s}
}
