package zmq

import (
	"context"
	"fmt"
	"github.com/go-zeromq/zmq4"
	"k8s.io/klog/v2"
)

func NewHub() *Hub {
	var hub = Hub{
		connections: make(map[string]*zmq4.Socket),
		commands:    make(chan *Command, 256),
		events:      make(chan *Event, 256),
		input:       make(chan *Message, 256),
		output:      make(chan *Message, 256),
	}
	go hub.commandProcessing()
	return &hub
}

type Hub struct {
	connections map[string]*zmq4.Socket
	commands    chan *Command
	events      chan *Event
	input       chan *Message
	output      chan *Message
}

func (hub Hub) Connections() map[string]*zmq4.Socket {
	return hub.connections
}

func (hub Hub) Connected(endpoint string) bool {
	_, has := hub.connections[endpoint]
	return has
}

func (hub Hub) Events() chan *Event {
	return hub.events
}

func (hub Hub) Commands() chan *Command {
	return hub.commands
}

func (hub Hub) Input() chan *Message {
	return hub.input
}

func (hub Hub) Output() chan *Message {
	return hub.output
}

func (hub Hub) ConnectedList() []string {
	keys := make([]string, len(hub.connections))
	i := 0
	for k := range hub.connections {
		keys[i] = k
		i++
	}
	return keys
}

func (hub Hub) tryDisconnect(endpoint string) {
	connection := hub.connections[endpoint] //todo добавление подключения
	err := (*connection).Close()
	if err != nil {
		//todo обработка ошибки
		return
	} else {
		hub.Events() <- &Event{Endpoint: endpoint, EventType: OnDisconnected}
		delete(hub.connections, endpoint)
	}
}

func (hub Hub) tryConnect(endpoint string) {
	i := 1

	if _, exists := hub.connections[endpoint]; exists {
		klog.Infof("Connection already exists in pool")
	} else {

		id := zmq4.SocketIdentity(fmt.Sprintf("dealer-%d", i))
		dealer := zmq4.NewDealer(context.Background(), zmq4.WithID(id))
		err := dealer.Dial(endpoint)
		if err != nil {
			// todo обработка ошибки
			klog.Error("dealer %d failed to dial: %v", i, err)
		} else {
			hub.connections[endpoint] = &dealer
			go hub.inputProcessing(endpoint) // todo контекст для выхода
			go hub.outputProcessing()        // todo контекст для выхода
			hub.events <- &Event{Endpoint: endpoint, EventType: OnConnected}
		}
	}

}

func (hub Hub) commandProcessing() {
	for {
		command := <-hub.commands
		switch t := command.CommandType; t {
		case TryConnect:
			hub.tryConnect(command.Endpoint)
		case TryDisconnect:
			hub.tryDisconnect(command.Endpoint)
		default:
			klog.Error(" ZMQ HUB command not supported", command.CommandType)
		}
	}
}

func (hub Hub) outputProcessing() {
	for {
		message := <-hub.output
		msg := zmq4.NewMsgFrom(message.Message)
		con := hub.connections[message.Endpoint] //todo обработать ошибку
		(*con).Send(msg)
	}
}

func (hub Hub) inputProcessing(endpoint string) { // @todo сделать контекст.
	for {
		con := hub.connections[endpoint]
		request, err := (*con).Recv()
		if err != nil {
			delete(hub.connections, endpoint)
			hub.events <- &Event{Endpoint: endpoint, EventType: OnDisconnected}
			klog.Error(err)
			break
		}
		message := request.Frames[0]
		hub.input <- &Message{message, endpoint}
	}
}
