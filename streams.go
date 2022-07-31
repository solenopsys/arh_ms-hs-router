package main

import (
	"errors"
	"k8s.io/klog/v2"
	"math/rand"
	"strconv"
)

type RoutingMap struct {
	wsConnections         map[uint32]string
	zmqConnections        map[uint32]string
	zmqConnectionsByTypes map[uint16][]string
}

func newRouting() *RoutingMap {
	return &RoutingMap{
		wsConnections:         map[uint32]string{},
		zmqConnections:        map[uint32]string{},
		zmqConnectionsByTypes: map[uint16][]string{},
	}
}

func getRandomServiceConnectionByType(connectionsMap map[uint16][]string, serviceType uint16, state *State) (string, error) {
	services := connectionsMap[serviceType]
	size := len(services)
	if size > 0 {
		return services[rand.Intn(size)], nil
	} else {
		return "", errors.New("not found")
	}

}

func (routingMap *RoutingMap) wsMessageProcessing(messageWrapper *MessageWrapper, state *State) {
	klog.Infof("PACK MESSAGE")
	klog.Infof(string(messageWrapper.body))
	var id = messageWrapper.meta.streamId

	if _, has := routingMap.wsConnections[id]; !has { //todo возврать ошибки если подключено к другому
		routingMap.wsConnections[id] = messageWrapper.key
	}

	if conn, has := routingMap.zmqConnections[id]; has { //todo возврать ошибки если подключено к другому
		state.messages <- &MessageWrapper{body: messageWrapper.body, messageType: ZmqOutput, key: conn, meta: messageWrapper.meta}
	} else {
		klog.Infof("SERVICE ID: " + strconv.Itoa(int(messageWrapper.meta.serviceId)))
		serviceType := messageWrapper.meta.serviceId
		randomZmqConnection, err := getRandomServiceConnectionByType(routingMap.zmqConnectionsByTypes, serviceType, state)
		if err == nil {
			routingMap.zmqConnections[id] = randomZmqConnection
			state.messages <- &MessageWrapper{body: messageWrapper.body, messageType: ZmqOutput, key: randomZmqConnection, meta: messageWrapper.meta}
		} else {
			body := append(messageWrapper.body[:8], []byte("ENDPOINT_NOT_FOUND")...)
			state.messages <- &MessageWrapper{body: body, messageType: WsOutput, key: messageWrapper.key, meta: messageWrapper.meta}
		}
	}
}

func (routingMap *RoutingMap) zqmMessageProcessing(messageWrapper *MessageWrapper, state *State) {
	connection := routingMap.wsConnections[messageWrapper.meta.streamId]
	state.messages <- &MessageWrapper{body: messageWrapper.body, messageType: WsOutput, key: connection, meta: messageWrapper.meta}
}

func (routingMap *RoutingMap) removeStreamsByWs(wsEndpoint string) {
	streamsIds := []uint32{}
	for id, endpoint := range routingMap.wsConnections {
		if endpoint == wsEndpoint {
			streamsIds = append(streamsIds, id)
		}
	}

	routingMap.removeStreamsByIds(streamsIds)
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func (routingMap *RoutingMap) sendCloseEvent(streamsIds []uint32, state *State) {
	for _, id := range streamsIds {
		klog.Infof("CLOSE EVENT SEND")
		connection := routingMap.wsConnections[id]

		state.incStat("EndpointDisconnected")

		// todo реализовать body := append(messageWrapper.body[:8], []byte("ENDPOINT_NOT_FOUND")...)
		state.messages <- &MessageWrapper{body: []byte("ENDPOINT_CLOSED_CONNECTION"), messageType: WsOutput, key: connection}
	}
}

func (routingMap *RoutingMap) removeStreamsByZmq(zmeEndpoint string, state *State) {
	streamsIds := []uint32{}
	for id, endpoint := range routingMap.zmqConnections {
		if endpoint == zmeEndpoint {
			streamsIds = append(streamsIds, id)
		}
	}

	routingMap.sendCloseEvent(streamsIds, state)
	routingMap.removeStreamsByIds(streamsIds)
	routingMap.removeEndpoint(zmeEndpoint)
}

func (routingMap *RoutingMap) removeEndpoint(zmeEndpoint string) {
	var found = false
	for key, endpoints := range routingMap.zmqConnectionsByTypes {
		for i, endpoint := range endpoints {
			if endpoint == zmeEndpoint {
				routingMap.zmqConnectionsByTypes[key] = removeIndex(endpoints, i)
				found = true
				break
			}
		}
		if found {
			break
		}
	}
}

func (routingMap *RoutingMap) addConnection(endpoint string, state *State) {
	typeId := state.endpoints[endpoint]
	if routingMap.zmqConnectionsByTypes[typeId] == nil {
		var intArray []string
		routingMap.zmqConnectionsByTypes[typeId] = intArray
	}
	current := routingMap.zmqConnectionsByTypes[typeId]
	item := append(current, endpoint)
	routingMap.zmqConnectionsByTypes[typeId] = item
}

func (routingMap *RoutingMap) removeStreamsByIds(streamsIds []uint32) {
	for _, id := range streamsIds {
		delete(routingMap.zmqConnections, id)
		delete(routingMap.wsConnections, id)
	}
}
