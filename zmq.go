package main

import (
	"context"
	"fmt"
	"github.com/go-zeromq/zmq4"
	"log"
	"time"
)

// ok
type ZmqHub struct {
	connected map[string]*zmq4.Socket
}

// ok
func newZqmHub() *ZmqHub {
	return &ZmqHub{
		connected: make(map[string]*zmq4.Socket),
	}
}

// ok
func (hub *ZmqHub) updateLoop(state *State, duration time.Duration) {
	for range time.Tick(duration) {
		// fmt.Println(now)
		id := "ZmqUpdateServices"
		state.command <- &Command{key: id, messageType: ZmqUpdateServices}
	}
}

// ok
func (hub *ZmqHub) unregisterDiff(state *State) {
	for key, _ := range hub.connected {
		if _, has := state.endpoints[key]; !has {
			state.command <- &Command{key: key, messageType: ZmqTryDisconnect}
		}
	}
}

// ok
func (hub *ZmqHub) registerDiff(state *State) {
	for endpoint, _ := range state.endpoints {
		if _, has := hub.connected[endpoint]; !has {
			state.command <- &Command{key: endpoint, messageType: ZmqTryConnect}
		}
	}
}

// ok
func (hub *ZmqHub) syncEndpoints(state *State) {
	hub.registerDiff(state)
	hub.unregisterDiff(state)
}

// ok
func (hub *ZmqHub) tryDisconnectProcessing(state *State, endpoint string) {
	connection := hub.connected[endpoint]
	err := (*connection).Close()
	if err != nil {
		return //todo обработка ошибки
	} else {
		state.command <- &Command{key: endpoint, messageType: ZmqOnDisconnected}
		delete(hub.connected, endpoint)
	}
}

// ok
func (hub *ZmqHub) tryConnectProcessing(state *State, endpoint string) {
	i := 1
	id := zmq4.SocketIdentity(fmt.Sprintf("dealer-%d", i))
	dealer := zmq4.NewDealer(context.Background(), zmq4.WithID(id))
	err := dealer.Dial(endpoint)
	if err != nil {
		log.Print("dealer %d failed to dial: %v", i, err) // todo обработка ошибки
	} else {
		hub.connected[endpoint] = &dealer
		state.command <- &Command{key: endpoint, messageType: ZmqOnConnected}
	}

}

//ok
func (hub *ZmqHub) onConnectProcessing(state *State, endpoint string) {

	//todo сюда
	go hub.zmqInMessageProcessing(state, endpoint)
	state.command <- &Command{key: endpoint, messageType: ZmqOnInitialized}
}

// ok
func (hub *ZmqHub) sendToZmqMessageProcessing(endpoint string, binaryMessage []byte) {
	msg := zmq4.NewMsgFrom(binaryMessage)
	con := hub.connected[endpoint]
	(*con).Send(msg)
}

//ok
func (hub *ZmqHub) zmqInMessageProcessing(state *State, endpoint string) { // @todo сделать контекст.
	for {
		con := hub.connected[endpoint]
		request, err := (*con).Recv()
		if err != nil {
			log.Print(err)
			state.command <- &Command{key: endpoint, messageType: ZmqOnDisconnected}
			break
		}
		message := request.Frames[0]
		state.messages <- &MessageWrapper{body: message, messageType: ZmqInput, key: endpoint}
	}
}
