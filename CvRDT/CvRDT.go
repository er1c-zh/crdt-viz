package CvRDT

type CvRDT interface {
	Query() interface{}
	Update() interface{}
	Merge() interface{}
}

type State interface {
	Init() // 初始化
}

type Inst struct {
	State State
	CvRDT
}

