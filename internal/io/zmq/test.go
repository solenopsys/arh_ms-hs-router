package zmq

import "github.com/go-zeromq/zmq4"

func NewZqmHub() *Hub {
	z := &Hub{
		connections: make(map[string]*zmq4.Socket),
		commands:    make(chan *Command),
		events:      make(chan *Event),
		input:       make(chan *Message),
		output:      make(chan *Message),
	}
	return z
}

func test() {
	var hub IO
	hub = NewZqmHub()

	println(hub)
}
