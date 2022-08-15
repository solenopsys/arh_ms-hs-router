package core

import (
	"errors"
	"k8s.io/utils/strings/slices"
	"math/rand"
	"solenopsys.org/zmq_router/internal/io/kube"
	"solenopsys.org/zmq_router/internal/io/zmq"
	"solenopsys.org/zmq_router/pkg/utils"
	"time"
)

type ServicesConf struct {
	services  map[string]uint16
	endpoints map[string]uint16
	groups    map[uint16][]string
}

type ServicesController struct {
	conf         *ServicesConf
	endpointsApi kube.EndpointsIntf
	mappingApi   kube.MappingIntf
	zmqHub       zmq.IO
}

func (sc ServicesController) UpdateEndpointsLoop(duration time.Duration) {
	for range time.Tick(duration) {
		sc.SyncEndpoints()
	}
}

func (sc ServicesController) unregisterDiff() {
	for _, endpoint := range sc.zmqHub.ConnectedList() {
		if _, has := sc.conf.endpoints[endpoint]; !has {
			sc.zmqHub.Commands() <- &zmq.Command{Endpoint: endpoint, CommandType: zmq.TryDisconnect}
			sc.removeEndpointFromGroup(endpoint) //todo возможно нужно обрабатывать событие
		}
	}
}

func (sc ServicesController) registerDiff() {
	for endpoint, _ := range sc.conf.endpoints {
		if connected := sc.zmqHub.Connected(endpoint); !connected {
			sc.zmqHub.Commands() <- &zmq.Command{Endpoint: endpoint, CommandType: zmq.TryConnect}
			sc.addEndpointToGroup(endpoint) //todo возможно нужно обрабатывать событие
		}
	}
}

func (sc ServicesController) SyncEndpoints() {

	endpoints := sc.endpointsApi.Endpoints()
	sc.conf.services = sc.mappingApi.Mapping()
	sc.conf.endpoints = convertEndpoints(endpoints, sc.conf.services)
	sc.registerDiff()
	sc.unregisterDiff()
}

func (sc ServicesController) getRandomServiceConnectionByType(serviceType uint16) (string, error) {
	services := sc.conf.groups[serviceType]
	size := len(services)
	if size > 0 {
		return services[rand.Intn(size)], nil
	} else {
		return "", errors.New("not found")
	}
}

func (sc ServicesController) removeEndpointFromGroup(zmeEndpoint string) {
	for key, endpoints := range sc.conf.groups {
		for i, endpoint := range endpoints {
			if endpoint == zmeEndpoint {
				sc.conf.groups[key] = utils.RemoveIndex(endpoints, i)
				return
			}
		}
	}
}

func (sc ServicesController) addEndpointToGroup(endpoint string) {
	typeId := sc.conf.endpoints[endpoint]
	if sc.conf.groups[typeId] == nil {
		sc.conf.groups[typeId] = []string{}
	}
	group := sc.conf.groups[typeId]
	if !slices.Contains(group, endpoint) {
		sc.conf.groups[typeId] = append(group, endpoint)
	}
}

func convertEndpoints(endpointsKeys map[string]string, services map[string]uint16) map[string]uint16 {
	res := make(map[string]uint16)
	for key, serviceName := range endpointsKeys {
		res[key] = services[serviceName]
	}
	return res
}
