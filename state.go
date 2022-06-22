package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
)

type CommandType int8

const (
	ZmqUpdateServices CommandType = iota
	ZmqTryConnect
	ZmqTryDisconnect
	ZmqOnConnected
	ZmqOnInitialized  // ? нет обработчика
	ZmqOnDisconnected // ? нет обработчика
	WsOnConnected
	WsOnInitialized // ? нет обработчика
	WsUpgradeError  // ? нет обработчика
	WsTryDisconnect
	WsOnDisconnected // ? нет обработчика
)

func (me CommandType) String() string {
	return [...]string{
		"Zmq UpdateServices",
		"Zmq TryConnect",
		"Zmq TryDisconnect",
		"Zmq OnConnected",
		"Zmq OnInitialized",
		"Zmq OnDisconnected",
		"Ws OnConnected",
		"Ws OnInitialized",
		"Ws UpgradeError",
		"Ws TryDisconnect",
		"Ws OnDisconnected",
	}[me]
}

type MessageType int8

const (
	WsInput MessageType = iota
	WsInputParsed
	ZmqInput
	ZmqInputParsed
	WsOutput
	ZmqOutput
)

func (me MessageType) String() string {
	return [...]string{
		"Ws Input",
		"Ws InputParsed",
		"Zmq Input",
		"Zmq InputParsed",
		"Ws Output",
		"Zmq Output",
	}[me]
}

type MessagePackage struct {
	streamId   uint32
	serviceId  uint16
	functionId uint16
}

type Command struct {
	messageType CommandType
	key         string
}

type MessageWrapper struct {
	body        []byte
	messageType MessageType
	userId      uint16
	meta        MessagePackage
	key         string
}

type CommandProcessingFunc func(command Command)
type MessageProcessingFunc func(message MessageWrapper)

type State struct {
	command      chan *Command
	messages     chan *MessageWrapper
	services     map[string]uint16
	endpoints    map[string]uint16
	commandsFunc map[CommandType][]CommandProcessingFunc
	messageFunc  map[MessageType]MessageProcessingFunc
	autoNext     map[MessageType]MessageType
}

func printMap(s string, m map[string]uint16) {
	bs, _ := json.Marshal(m)
	fmt.Println()

	println(s, string(bs))
}

func newState(zmq *ZmqHub, ws *WsHub, routing *RoutingMap) *State {
	state := State{
		command:      make(chan *Command, 256),
		messages:     make(chan *MessageWrapper, 256),
		services:     make(map[string]uint16),
		endpoints:    make(map[string]uint16),
		commandsFunc: make(map[CommandType][]CommandProcessingFunc),
		messageFunc:  make(map[MessageType]MessageProcessingFunc),
		autoNext:     make(map[MessageType]MessageType),
	}

	state.autoNext[WsInput] = WsInputParsed
	state.autoNext[ZmqInput] = ZmqInputParsed

	state.commandsFunc[ZmqUpdateServices] = []CommandProcessingFunc{
		func(command Command) {
			var services map[string]uint16
			var endpoints map[string]uint16
			if devMode {
				services, endpoints = getEndpointsDeveloperMode()
			} else {
				services, endpoints = getEndpoints()
			}
			printMap("SERVICES ", services)
			state.services = services
			printMap("ENDPOINTS ", endpoints)
			state.endpoints = endpoints

			zmq.syncEndpoints(&state)

		}}
	state.commandsFunc[ZmqTryConnect] = []CommandProcessingFunc{
		func(command Command) { zmq.tryConnectProcessing(&state, command.key) }}
	state.commandsFunc[ZmqTryDisconnect] = []CommandProcessingFunc{
		func(command Command) { zmq.tryDisconnectProcessing(&state, command.key) }}
	state.commandsFunc[ZmqOnConnected] = []CommandProcessingFunc{
		func(command Command) {
			zmq.onConnectProcessing(&state, command.key)
			routing.addConnection(command.key, &state)
		}}
	state.commandsFunc[ZmqOnDisconnected] = []CommandProcessingFunc{
		func(command Command) { routing.removeStreamsByZmq(command.key, &state) },
		func(command Command) { delete(zmq.connected, command.key) },
	}

	state.commandsFunc[WsOnConnected] = []CommandProcessingFunc{
		func(command Command) { ws.onConnectProcessing(&state, command.key) }}
	state.commandsFunc[WsTryDisconnect] = []CommandProcessingFunc{
		func(command Command) { ws.tryDisconnectProcessing(&state, command.key) }}
	state.commandsFunc[WsOnDisconnected] = []CommandProcessingFunc{
		func(command Command) { routing.removeStreamsByWs(command.key) },
		func(command Command) { delete(ws.clients, command.key) },
	}

	state.messageFunc[WsInput] = func(message MessageWrapper) { state.extractMeta(&message) }
	state.messageFunc[WsInputParsed] = func(message MessageWrapper) { routing.wsMessageProcessing(&message, &state) }
	state.messageFunc[WsOutput] = func(message MessageWrapper) { ws.sendToWsMessageProcessing(message.key, message.body) }

	state.messageFunc[ZmqInput] = func(message MessageWrapper) { state.extractMeta(&message) }
	state.messageFunc[ZmqInputParsed] = func(message MessageWrapper) { routing.zqmMessageProcessing(&message, &state) }
	state.messageFunc[ZmqOutput] = func(message MessageWrapper) { zmq.sendToZmqMessageProcessing(message.key, message.body) }

	return &state
}

func (s *State) commandProcessor() {
	for {
		command := <-s.command

		log.Println("COMMAND:", command.messageType, command.key)

		if commandFuncs, ok := s.commandsFunc[command.messageType]; ok {
			for _, commandFunc := range commandFuncs {
				go commandFunc(*command)
			}
		}
	}
}

func (s *State) messageProcessor() {
	for {
		message := <-s.messages
		log.Println("MESSAGE:", message.messageType, message.key)
		if messageFunc, ok := s.messageFunc[message.messageType]; ok {
			messageFunc(*message)
		}
	}
}

func (state *State) extractMeta(message *MessageWrapper) {
	streamId := binary.BigEndian.Uint32(message.body[:4])
	serviceId := binary.BigEndian.Uint16(message.body[4:6])
	functionId := binary.BigEndian.Uint16(message.body[6:8])
	meta := MessagePackage{streamId: streamId, serviceId: serviceId, functionId: functionId}
	state.messages <- &MessageWrapper{body: message.body, messageType: state.autoNext[message.messageType], meta: meta, key: message.key}
}
