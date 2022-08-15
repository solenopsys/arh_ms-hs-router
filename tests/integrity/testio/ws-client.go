package testio

import (
	"github.com/gorilla/websocket"
	"log"
)

type WsStatus int

const (
	New WsStatus = iota
	Connected
	Error
	Closed
)

type WsMessage struct {
	ClientId uint16
	Body     []byte
}

type WsClient struct {
	Connection *websocket.Conn
	ClientId   uint16
	State      WsStatus
	Err        error
}

type WsClientsPool struct {
	ToWs              chan *WsMessage
	FromWs            chan *WsMessage
	ClientConnections map[uint16]*WsClient
}

func (client *WsClient) tryConnection(url string, clientId uint16) error {
	var err error
	client.Connection, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		client.State = Error
		client.Err = err
		log.Print(err)
	} else {
		log.Print("Ws succeful Connection ", url)
		log.Print("ST1 ", client.State)
		client.State = Connected
		log.Print("ST2 ", client.State)
	}
	return err
}

func (client WsClient) disconnect() {
	err := client.Connection.Close()

	if err != nil {
		client.State = Error
		client.Err = err
		log.Panic(err)
	}
}

func (client WsClient) listen(pipe chan *WsMessage) {
	for {
		_, message, err := client.Connection.ReadMessage()
		if err != nil {

			client.State = Error
		} else {
			pipe <- &WsMessage{client.ClientId, message}
		}
	}
}

func (client WsClient) sendMessage(message []byte) {
	log.Print("STATE", client.State)
	if client.State == Connected {
		log.Print("MESSAGE WRITE")
		err := client.Connection.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.State = Error
			client.Err = err
			log.Panic(err)
		}
	}
}

func (pool WsClientsPool) CreateClient(url string, clientId uint16) (*WsClient, error) {
	client := NewClient()
	pool.ClientConnections[clientId] = client
	err := client.tryConnection(url, clientId)

	if err == nil {
		go client.listen(pool.FromWs)
	}

	return client, err
}

func (pool WsClientsPool) deleteClient(clientId uint16) {
	if client, ok := pool.ClientConnections[clientId]; ok {
		if client.State == Closed {
			delete(pool.ClientConnections, clientId)
		}
	}
}

func (pool WsClientsPool) sendMessage(clientId uint16, message []byte) {
	if client, ok := pool.ClientConnections[clientId]; ok {
		client.sendMessage(message)
	}
}

func (pool WsClientsPool) sendLoop() {
	for {
		message := <-pool.ToWs
		println("TO MESSAGE")
		pool.sendMessage(message.ClientId, message.Body)
	}
}

func NewClient() *WsClient {
	return &WsClient{State: New}
}

func NewTestWsPool() *WsClientsPool {
	pool := WsClientsPool{
		FromWs:            make(chan *WsMessage, 256),
		ToWs:              make(chan *WsMessage, 256),
		ClientConnections: make(map[uint16]*WsClient),
	}
	go pool.sendLoop()
	return &pool
}
