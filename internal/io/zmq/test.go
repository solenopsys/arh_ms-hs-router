package zmq

import "github.com/go-zeromq/zmq4"

func NewZqmHub() *Hub {
	z := &Hub{
		Connections: make(map[string]*zmq4.Socket),
		Commands:    make(chan *Command),
		Events:      make(chan *Event),
		Input:       make(chan *Message),
		Output:      make(chan *Message),
	}
	return z
}

func test() {
	var hub IO
	hub = NewZqmHub()

	println(hub)
}
