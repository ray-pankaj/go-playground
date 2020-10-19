package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pingPeriod = 5 * time.Second
)

type peer struct {
	ws       *websocket.Conn
	send     chan *Message
	username string
}

func (p *peer) writer() {
	pingTicker := time.NewTicker(pingPeriod)
	ws := p.ws
	defer func() {
		pingTicker.Stop()
		// TODO: Should both reader and writer close sockets? Closing an already closed socket?
		ws.Close()
	}()
	for {
		select {
		case <-pingTicker.C:
			if err := ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"message_type":"SERVER_PING","data":"%v"}`, time.Now()))); err != nil {
				return
			}
		case msg := <-p.send:
			ws.WriteJSON(msg)
		}
	}
}

func (p *peer) reader() {
	ws := p.ws
	defer ws.Close()
	ws.SetReadLimit(512)
	var msg Message
	for {
		_, body, err := ws.ReadMessage()
		log.Println(body, err)
		if err != nil {
			log.Println(err)
		}
		if err := json.Unmarshal(body, &msg); err != nil {
			log.Println(err)
		} else {
			log.Println(body)
			log.Println((fmt.Sprintf("%T %v", msg.Data, msg.Data)))
			log.Println(msg.MessageType)
			log.Println(msg.Room)
			h.messageRouter <- &MessageWrapper{&msg, p}
		}
	}
}
