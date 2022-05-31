package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"encoding/json"
)

type Action func(s *Server)

func Broadcast(msg string) Action {
	return func(s *Server) {
		for ws := range s.conns {
			if err := websocket.Message.Send(ws, msg); err != nil {
				fmt.Println("Can't send")
			}
		}
	}
}

func AddConn(ws *websocket.Conn) Action {
	return func(s *Server) {
		s.conns[ws] = struct{}{}
		go func() {
			userConnectedMessage := &Message{
				Type: "text",
				Data: "New user connected",
			}
			b, err := json.Marshal(userConnectedMessage)
			if err != nil {
				fmt.Println(err)
				return
			}
			s.dispatch <- Broadcast(string(b))
		}()
	}
}

func RemoveConn(ws *websocket.Conn) Action {
	return func(s *Server) {
		delete(s.conns, ws)
		go func() {
			userDisconnectedMessage := &Message{
				Type: "text",
				Data: "User disconnected",
			}
			b, err := json.Marshal(userDisconnectedMessage)
			if err != nil {
				fmt.Println(err)
				return
			}
			s.dispatch <- Broadcast(string(b))
		}()
	}
}