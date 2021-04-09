package view

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Logger interface {
	Info(fmt string, args... interface{})
}

type Server struct{
	Logger
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
	if err != nil {
		s.Info("upgrade fail: %s", err.Error())
		return
	}
	for {

		conn.ReadJSON()
	}
}
