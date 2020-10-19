package main

type subscription struct {
	conn *peer
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]*room

	// Inbound messages from the connections.

	// Register requests from the connections.
	addRoom chan *room

	// Unregister requests from connections.
	delRoom chan *room

	messageRouter chan *MessageWrapper
}

var h = hub{
	rooms:         make(map[string]*room),
	addRoom:       make(chan *room),
	delRoom:       make(chan *room),
	messageRouter: make(chan *MessageWrapper),
}

func (h *hub) run() {
	for {
		select {
		case r := <-h.addRoom:
			// Possible Race condition?
			h.rooms[r.Name] = r
		case r := <-h.delRoom:
			delete(h.rooms, r.Name)

		case mw := <-h.messageRouter:
			msg := mw.msg
			p := mw.p
			switch msg.MessageType {
			case "CREATE_ROOM":
				h.rooms[msg.Room] = CreateRoom(msg.Room)
			case "JOIN_ROOM":
				h.rooms[msg.Room].JoinChannel <- p
			case "CHAT":
				//TODO: Check websocket.PreparedMessage for broadcast
				h.rooms[msg.Room].SendChannel <- msg
			case "WEBRTC_OFFER":
				h.rooms[msg.Room].SendChannel <- msg
			case "WEBRTC_ANSWER":
				h.rooms[msg.Room].SendChannel <- msg
			}
		}
	}

}
