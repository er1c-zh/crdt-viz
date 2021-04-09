package view

const (
	ActionInit = 1
)

type Request struct {
	Action int `json:"action"`
}

type Node struct {
	Addr string `json:"addr"`
	Desc string `json:"desc"`
}

type Msg struct {
	Desc string `json:"desc"`
	From string `json:"from"`
	To   string `json:"to"`
}

type Graph struct {
	NodeList []Node `json:"node_list"`
	MsgList  []Msg  `json:"msg_list"`
}
