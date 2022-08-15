package testio

import (
	"context"
	"github.com/go-zeromq/zmq4"
	"log"
	"strconv"
	"sync"
)

type ZmqStatus int

const (
	ZmqNew ZmqStatus = iota
	ZmqListened
	ZmqError
	ZmqClosed
)

type ZmqMessage struct {
	Port    uint16
	Address []byte
	Message []byte
}

type ZmqServer struct {
	Port   uint16
	State  ZmqStatus
	Err    error
	Wg     sync.WaitGroup
	router zmq4.Socket
}

type ZmqServersPool struct {
	ToHub   chan *ZmqMessage
	FromHub chan *ZmqMessage
	Servers map[uint16]*ZmqServer
}

func newZqmServer() *ZmqServer {
	var wg sync.WaitGroup
	return &ZmqServer{Wg: wg, State: ZmqNew}
}

func (server *ZmqServer) openPort(url string, port uint16) {
	//server.Wg.Add(1)
	server.Port = port
	startFunc := func() {
		log.Print("START SERVER")
		ps := strconv.FormatInt(int64(port), 10)
		server.router = zmq4.NewRouter(context.Background(), zmq4.WithID(zmq4.SocketIdentity("router"+ps)))
		servAddress := url + ":" + ps
		err := server.router.Listen(servAddress)
		if err != nil {

			server.State = ZmqError
			server.Err = err

			log.Panic("OPER PORT ERROR SERVER", err)
		} else {
			log.Println("Server started " + servAddress)
		}
		server.State = ZmqListened
		//	server.Wg.Wait()
	}

	startFunc()
}

func (server *ZmqServer) close() {
	server.Wg.Done()
	server.State = ZmqClosed
}

func (server *ZmqServer) listen(pipe chan *ZmqMessage) {
	for true {
		request, err := (server.router).Recv()
		if err != nil {
			server.State = ZmqError
			server.Err = err
			server.Wg.Done()
		}
		pipe <- &ZmqMessage{
			server.Port,
			request.Frames[0],
			request.Frames[1],
		}
	}
}

func (server *ZmqServer) send(message *ZmqMessage) {
	msg := zmq4.NewMsgFrom(message.Address, message.Message)
	err := (server.router).Send(msg)
	if err != nil {
		server.Err = err
		server.Wg.Done()
	}
}

func (pool *ZmqServersPool) CreateServer(url string, port uint16) *ZmqServer {
	server := newZqmServer()
	server.openPort(url, port)
	go server.listen(pool.FromHub)

	return server
}

func (pool *ZmqServersPool) deleteServer(port uint16) {
	if server, ok := pool.Servers[port]; ok {
		if server.State == ZmqClosed {
			delete(pool.Servers, port)
		}
	}
}

func (pool *ZmqServersPool) sendMessage(message *ZmqMessage) {
	if server, ok := pool.Servers[message.Port]; ok {
		server.send(message)
	}
}

func (pool *ZmqServersPool) sendLoop() {
	for {
		message := <-pool.ToHub
		println("TO HUB MESSAGE")
		pool.sendMessage(message)
	}
}

func NewTestZmqPool() *ZmqServersPool {
	return &ZmqServersPool{
		FromHub: make(chan *ZmqMessage, 256),
		ToHub:   make(chan *ZmqMessage, 256),
		Servers: make(map[uint16]*ZmqServer),
	}
}
