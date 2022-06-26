package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
)

type Action func(s *Server)

type ConnSet map[*websocket.Conn]struct{}

type Server struct {
	conns    ConnSet
	dispatch chan Action
}

func broadcast(message string) Action {
	return func(s *Server) {
		for ws := range s.conns {
			if err := websocket.Message.Send(ws, message); err != nil {
				fmt.Println(err, "Can't send message")
			}
		}
	}
}

func sendUserCount(s *Server) {
	s.dispatch <- broadcast(`{"type" : "userCount", "data" : ` + strconv.Itoa(len(s.conns)) + `}`)
}

func addConn(ws *websocket.Conn) Action {
	return func(s *Server) {
		s.conns[ws] = struct{}{}
		go sendUserCount(s)
	}
}

func removeConn(ws *websocket.Conn) Action {
	return func(s *Server) {
		delete(s.conns, ws)
		go sendUserCount(s)
	}
}

func (s *Server) handler(ws *websocket.Conn) {
	s.dispatch <- addConn(ws)

	defer func() { s.dispatch <- removeConn(ws) }()

	for {
		var reply string

		if err := websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println(err, "Can't receive message")
			break
		}
		fmt.Println(reply)
		s.dispatch <- broadcast(reply)
	}
}

func (s *Server) runDispatcher() {
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

	go s.runDispatcher()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/ws", websocket.Handler(s.handler))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
