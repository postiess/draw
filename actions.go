package main

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type Action func(s *Server)

func Broadcast(str string) Action {
	return func(s *Server) {
		for ws := range s.conns {
			if err := websocket.Message.Send(ws, str); err != nil {
				fmt.Println("Can't send")
			}
		}
	}
}

func AddConn(ws *websocket.Conn) Action {
	return func(s *Server) {
		s.conns[ws] = struct{}{}
		go func() {
			s.dispatch <- Broadcast("New user connected")
		}()
	}
}

func RemoveConn(ws *websocket.Conn) Action {
	return func(s *Server) {
		delete(s.conns, ws)
		go func() {
			s.dispatch <- Broadcast("User disconnected")
		}()
	}
}