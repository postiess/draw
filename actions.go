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

func sendUserCount(s *Server){
	userMessage := &Message{
		Type: "userCount",
		Data: len(s.conns),
	}
	b, err := json.Marshal(userMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.dispatch <- Broadcast(string(b))
}

func AddConn(ws *websocket.Conn) Action {
	return func(s *Server) {
		s.conns[ws] = struct{}{}
		go sendUserCount(s)
	}
}

func RemoveConn(ws *websocket.Conn) Action {
	return func(s *Server) {
		delete(s.conns, ws)
		go sendUserCount(s)
	}
}