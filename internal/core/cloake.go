package core

//func sendCloseEvent(streamsIds []uint32) {
//	for _, id := range streamsIds {
//		klog.Infof("CLOSE EVENT SEND")
//		connection := model.WsStreamsConnections[id]
//
//		// todo реализовать body := append(messageWrapper.body[:8], []byte("ENDPOINT_NOT_FOUND")...)
//		state.Messages <- &MessageWrapper{Body: []byte("ENDPOINT_CLOSED_CONNECTION"), MessageType: WsOutput, ConnectionKey: connection}
//	}
//}

/*	state.Messages <- &MessageWrapper{
	Body:          messageWrapper.Body,
	MessageType:   WsOutput,
	ConnectionKey: connection,
	Meta:          messageWrapper.Meta}*/

//func (r Routing) sendMessage2(messageWrapper *MessageWrapper, conn StreamConfig) {
//	state.Messages <- &MessageWrapper{
//		Body:          messageWrapper.Body,
//		MessageType:   ZmqOutput,
//		ConnectionKey: conn,
//		Meta:          messageWrapper.Meta}
//}

/*
вынести в стейт
	m.sendCloseEvent(streamsIds)
	m.removeEndpointFromGroup(zmeEndpoint)
*/

//func (cd Coordinator) SyncConnectionScheduler() {
//	for range time.Tick(10) {
//		cd.statistic.Increment <- "SyncZmqConnections"
//		changed := cd.syncZmqConnections()
//
//		if changed {
//			cd.statistic.Increment <- "SyncZmqConnectionsChanged"
//		}
//	}
//}
//
//func (cd Coordinator) syncZmqConnections() bool {
//	var changed = false
//
//	for endpoint, _ := range cd.endpoints.EndpointsMap() {
//		if !cd.zmqHub.Connected(endpoint) {
//			changed = true
//			cd.zmqHub.GetCommands() <- &zmq.Command{Endpoint: endpoint, CommandType: zmq.TryConnect}
//		}
//	}
//
//	return changed
//	// todo disconnected
//}

/*
	klog.Infof("INIT MESSAGE")
	klog.Infof(string(messageWrapper.Body))
*/
//func (r Routing) sendMessage(messageWrapper *MessageWrapper, randomZmqConnection string) {
//	body := messageWrapper.Body
//	binary.BigEndian.PutUint16(body[6:], messageWrapper.UserId)
//	state.Messages <- &MessageWrapper{
//		Body:          body,
//		MessageType:   ZmqOutput,
//		ConnectionKey: randomZmqConnection,
//		Meta:          messageWrapper.Meta}
//}

//return &MessageWrapper{
//	Body:         ,
//	MessageType:   WsOutput,
//	ConnectionKey: message.ConnectionKey,
//	Meta:          message.Meta,
//}
//var next MessageType
//state.Messages <- &MessageWrapper{
//		Body:          message.Body,
//		MessageType:   next,
//		Meta:          meta,
//		ConnectionKey: message.ConnectionKey,
//	}*/

//func UpdateConfigMapLoop(state *core.State, duration time.Duration) {
//	for range time.Tick(duration) {
//		id := "UpdateConfigMap"
//		state.Command <- &core.Command{Key: id, MessageType: core.UpdateConfigMap}
//	}
//}

//func (s *State) commandProcessor() {
//	for {
//		command := <-s.Command
//
//		klog.Infof("COMMAND: %", command.Key)
//
//		if commandFuncs, ok := s.CommandsFunc[command.MessageType]; ok {
//			for _, commandFunc := range commandFuncs {
//				go commandFunc(*command)
//			}
//		}
//	}
//}

//func (s *State) IncStat(val string) {
//	s.StatPipe <- val
//}
//
//func (s *State) messageProcessor() {
//	for {
//		message := <-s.Messages
//		klog.Infof("MESSAGE:", message.MessageType, message.ConnectionKey)
//		if MessageFunc, ok := s.MessageFunc[message.MessageType]; ok {
//			MessageFunc(*message)
//		}
//	}
//}
//
//
//
//

//func(command Command) {
//	var endpointsKeys map[string]string
//	var err error
//	endpointsKeys, err = Apif.GetEndpoints()
//
//	var services map[string]uint16
//	services, err = Apif.GetHsMapping()
//	if err != nil {
//		klog.Infof("Get data error", err)
//	} else {
//		state.ServicesMap = services
//		endpoints := convertEndpoints(endpointsKeys, services)
//	}
//}

//func (r Routing) removeStreamsByZmq(zmeEndpoint string, message string) { //todo испраивить message на код
//	for id, streamConfig := range r.Streams {
//		if streamConfig.ZmqEndpoint == zmeEndpoint {
//			r.Remove <- &RemoveParam{id, message}
//		}
//	}
//}
//
//func (r Routing) removeStreamsByWs(wsEndpoint string, message string) { //todo испраивить message на код
//	for id, streamConfig := range r.Streams {
//		if streamConfig.WsEndpoint == wsEndpoint {
//			r.Remove <- &RemoveParam{id, message}
//		}
//	}
//}
//	sc.zmqHub.GetCommands() <- &zmq.Command{Endpoint: endpoint, CommandType: zmq.TryDisconnect}
//			sc.zmqHub.GetCommands() <- &zmq.Command{Endpoint: endpoint, CommandType: zmq.TryConnect}

/*BeforeEach(func() {
	endpoint = ""
	services =&core.ServicesController{
		EndpointsApi: EndpointsMock{},
		MappingApi:   MappingMock{},
		ServicesMap:  make(map[string]uint16),
		EndpointsMap: make(map[string]uint16),
		Groups:       make(map[uint16][]string),
		Rand:         func(size int) int { return rand.Intn(size) },
	}
})

func printZmqMessages() {
	for {
		message := <-model.zmqPool.FromHub

		println("FROM HUB-- -------", string(message.Message))
		println("PORT", strconv.FormatUint(uint64(message.Port), 10))
	}

}

*/
