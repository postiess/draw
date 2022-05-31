package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Action func(s *Server)

type ConnSet map[*websocket.Conn]struct{}

type Server struct {
	conns    ConnSet
	dispatch chan Action
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func Broadcast(message string) Action {
	return func(s *Server) {
		for ws := range s.conns {
			if err := websocket.Message.Send(ws, message); err != nil {
				fmt.Println("Can't send message")
			}
		}
	}
}

func sendUserCount(s *Server) {
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

func (s *Server) Handler(ws *websocket.Conn) {
	s.dispatch <- AddConn(ws)

	defer func() { s.dispatch <- RemoveConn(ws) }()

	for {
		var reply string

		if err := websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println(err)
			fmt.Println("Can't receive")
			break
		}
		fmt.Println(reply)
		s.dispatch <- Broadcast(reply)
	}
}

func (s *Server) RunDispatcher() {
	for {
		action := <-s.dispatch
		action(s)
	}
}

func main() {
	s := &Server{
		conns:    ConnSet{},
		dispatch: make(chan Action),
	}

	go s.RunDispatcher()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/ws", websocket.Handler(s.Handler))
	err := http.ListenAndServe(":3000", nil);
	if  err != nil {
		log.Fatal(err)
	}
}
