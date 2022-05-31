package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type ConnSet map[*websocket.Conn]struct{}

type Server struct {
	conns    ConnSet
	dispatch chan Action
}

type Message struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}

func (s *Server) Handler(ws *websocket.Conn) {

	var err error
	s.dispatch <- AddConn(ws)
	defer func() { s.dispatch <- RemoveConn(ws) }()

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
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

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.Handle("/ws", websocket.Handler(s.Handler))
	log.Print("Starting on :3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
