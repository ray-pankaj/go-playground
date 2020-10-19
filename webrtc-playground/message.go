package main

import (
	"encoding/json"
)

type Message struct {
	From        string          `json:"from"`
	To          string          `json:"to"`
	Data        json.RawMessage `json:"data"`
	MessageType string          `json:"message_type"`
	Room        string          `json:"room"`
}

type MessageWrapper struct {
	msg *Message
	p   *peer
}
