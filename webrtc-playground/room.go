package main

import "log"

type room struct {
	Name         string
	Peers        map[string]*peer
	JoinChannel  chan *peer
	LeaveChannel chan *peer
	SendChannel  chan *Message
}

func (r *room) init() {
	for {
		select {
		case p := <-r.JoinChannel:
			// Possible Race condition?
			r.Peers[p.username] = p
		case p := <-r.LeaveChannel:
			delete(r.Peers, p.username)
		case msg := <-r.SendChannel:
			switch msg.To {
			case "broadcast":
				for _, p := range r.Peers {
					if p.username != msg.From {
						p.send <- msg
					}
				}
			case "server":
				log.Println(msg)
			default:
				r.Peers[msg.To].send <- msg
			}
		}
	}
}

func CreateRoom(name string) *room {
	r := &room{
		Name:         name,
		Peers:        make(map[string]*peer),
		JoinChannel:  make(chan *peer),
		LeaveChannel: make(chan *peer),
		SendChannel:  make(chan *Message),
	}
	go r.init()
	return r
}
