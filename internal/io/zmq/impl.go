package zmq

import (
	"context"
	"fmt"
	"github.com/go-zeromq/zmq4"
	"k8s.io/klog/v2"
)

type Hub struct {
	Connections map[string]*zmq4.Socket
	Commands    chan *Command
	Events      chan *Event
	Input       chan *Message
	Output      chan *Message
}

func (hub Hub) GetConnections() map[string]*zmq4.Socket {
	return hub.Connections
}

func (hub Hub) Connected(endpoint string) bool {
	_, has := hub.Connections[endpoint]
	return has
}

func (hub Hub) GetEvents() chan *Event {
	return hub.Events
}

func (hub Hub) GetCommands() chan *Command {
	return hub.Commands
}

func (hub Hub) GetInput() chan *Message {
	return hub.Input
}

func (hub Hub) GetOutput() chan *Message {
	return hub.Output
}

func (hub Hub) ConnectedList() []string {
	keys := make([]string, len(hub.Connections))
	i := 0
	for k := range hub.Connections {
		keys[i] = k
		i++
	}
	return keys
}

func (hub Hub) tryDisconnect(endpoint string) {
	connection := hub.Connections[endpoint] //todo добавление подключения
	err := (*connection).Close()
	if err != nil {
		//todo обработка ошибки
		return
	} else {
		hub.GetEvents() <- &Event{Endpoint: endpoint, EventType: OnDisconnected}
		delete(hub.Connections, endpoint)
	}
}

func (hub Hub) tryConnect(endpoint string) {
	i := 1

	if _, exists := hub.Connections[endpoint]; exists {
		klog.Infof("Connection already exists in pool")
	} else {

		id := zmq4.SocketIdentity(fmt.Sprintf("dealer-%d", i))
		dealer := zmq4.NewDealer(context.Background(), zmq4.WithID(id))
		err := dealer.Dial(endpoint)
		if err != nil {
			// todo обработка ошибки
			klog.Error("dealer %d failed to dial: %v", i, err)
		} else {
			hub.Connections[endpoint] = &dealer
			go hub.inputProcessing(endpoint) // todo контекст для выхода
			go hub.outputProcessing()        // todo контекст для выхода
			hub.Events <- &Event{Endpoint: endpoint, EventType: OnConnected}
		}
	}

}

func (hub Hub) CommandProcessing() {
	for {
		command := <-hub.Commands
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
		message := <-hub.Output
		msg := zmq4.NewMsgFrom(message.Message)
		con := hub.Connections[message.Endpoint] //todo обработать ошибку
		(*con).Send(msg)
	}
}

func (hub Hub) inputProcessing(endpoint string) { // @todo сделать контекст.
	for {
		con := hub.Connections[endpoint]
		request, err := (*con).Recv()
		if err != nil {
			delete(hub.Connections, endpoint)
			hub.Events <- &Event{Endpoint: endpoint, EventType: OnDisconnected}
			klog.Error(err)
			break
		}
		message := request.Frames[0]
		hub.Input <- &Message{message, endpoint}
	}
}
