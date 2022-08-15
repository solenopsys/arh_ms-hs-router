package ws

import (
	"github.com/gorilla/websocket"
	"k8s.io/klog/v2"
)

type Connection struct {
	connection *websocket.Conn
	userId     uint16
}

type Hub struct {
	connections map[string]*Connection
	commands    chan *Command
	events      chan *Event
	input       chan *Message
	output      chan *Message
}

func NewHub() *Hub {
	var hub = Hub{
		connections: make(map[string]*Connection),
		commands:    make(chan *Command, 256),
		events:      make(chan *Event, 256),
		input:       make(chan *Message, 256),
		output:      make(chan *Message, 256),
	}
	go hub.sendToWsMessageProcessing()
	return &hub
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

func (wsHub Hub) AddConnection(connection *websocket.Conn, userId uint16) {
	var id = connection.RemoteAddr().String()
	wsHub.connections[id] = &Connection{userId: userId, connection: connection}

	go wsHub.WsInMessageProcessing(id)
	wsHub.events <- &Event{Key: id, EventType: OnConnected}
}

func (wsHub Hub) TryDisconnectProcessing(Key string) {
	if client, ok := wsHub.connections[Key]; ok {
		delete(wsHub.connections, Key)
		err := client.connection.Close()
		if err != nil {
			klog.Error("Disconnect error:", err)
		} else {
			wsHub.events <- &Event{Key: Key, EventType: OnDisconnected}
		}
	}
}

func (wsHub Hub) WsInMessageProcessing(key string) { //todo обработать отключение через контекст.
	conn := wsHub.connections[key]
	for {
		_, message, err := conn.connection.ReadMessage()
		if err != nil {
			klog.Error("Read message error:", err)
			wsHub.events <- &Event{Key: key, EventType: OnDisconnected}
			break
		}
		klog.Infof("Ws in мessage:", string(message))
		wsHub.input <- &Message{
			Message:       message,
			ConnectionKey: key,
			User:          conn.userId,
			From:          key,
		}
	}
}

func (wsHub Hub) sendToWsMessageProcessing() {
	for {
		message := <-wsHub.output
		wsHub.sendToWsMessage(message.ConnectionKey, message.Message)
	}
}

func (wsHub Hub) sendToWsMessage(Key string, body []byte) {
	conn := wsHub.connections[Key]
	if conn != nil {
		err := conn.connection.WriteMessage(websocket.BinaryMessage, body)
		if err != nil {
			klog.Error("read:", err)
		}
	}
}
