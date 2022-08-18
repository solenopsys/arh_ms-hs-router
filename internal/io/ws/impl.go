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
	Connections map[string]*Connection
	Commands    chan *Command
	Events      chan *Event
	Input       chan *Message
	Output      chan *Message
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

func (wsHub Hub) AddConnection(connection *websocket.Conn, userId uint16) {
	var id = connection.RemoteAddr().String()
	wsHub.Connections[id] = &Connection{userId: userId, connection: connection}

	go wsHub.WsInMessageProcessing(id)
	wsHub.Events <- &Event{Key: id, EventType: OnConnected}
}

func (wsHub Hub) TryDisconnectProcessing(Key string) {
	if client, ok := wsHub.Connections[Key]; ok {
		delete(wsHub.Connections, Key)
		err := client.connection.Close()
		if err != nil {
			klog.Error("Disconnect error:", err)
		} else {
			wsHub.Events <- &Event{Key: Key, EventType: OnDisconnected}
		}
	}
}

func (wsHub Hub) WsInMessageProcessing(key string) { //todo обработать отключение через контекст.
	conn := wsHub.Connections[key]
	for {
		_, message, err := conn.connection.ReadMessage()
		if err != nil {
			klog.Error("Read message error:", err)
			wsHub.Events <- &Event{Key: key, EventType: OnDisconnected}
			break
		}
		klog.Infof("Ws in мessage:", string(message))
		wsHub.Input <- &Message{
			Message:       message,
			ConnectionKey: key,
			User:          conn.userId,
			From:          key,
		}
	}
}

func (wsHub Hub) SendToWsMessageProcessing() {
	for {
		message := <-wsHub.Output
		wsHub.sendToWsMessage(message.ConnectionKey, message.Message)
	}
}

func (wsHub Hub) sendToWsMessage(Key string, body []byte) {
	conn := wsHub.Connections[Key]
	if conn != nil {
		err := conn.connection.WriteMessage(websocket.BinaryMessage, body)
		if err != nil {
			klog.Error("read:", err)
		}
	}
}
