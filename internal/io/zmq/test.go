package zmq

import "github.com/go-zeromq/zmq4"

func NewZqmHub() *Hub {
	z := &Hub{
		Connections: make(map[string]*zmq4.Socket),
		Commands:    make(chan *Command, 256),
		Events:      make(chan *Event, 256),
		Input:       make(chan *Message, 256),
		Output:      make(chan *Message, 256),
	}
	return z
}

func test() {
	var hub IO
	hub = NewZqmHub()

	println(hub)
}
