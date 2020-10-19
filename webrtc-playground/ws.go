package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	p := peer{ws, make(chan *Message), username}
	go p.writer()
	go p.reader()
}
