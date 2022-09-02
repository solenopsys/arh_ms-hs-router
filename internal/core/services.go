package core

import (
	"errors"
	"k8s.io/utils/strings/slices"
	"math/rand"
	"reflect"
	"solenopsys.org/zmq_router/internal/io/kube"
)

type Services interface {
	Services() map[string]uint16
	Endpoints() map[string]uint16
	RandomEndpointByType(serviceType uint16) (string, error)
	SyncEndpoints()
}

type ServicesController struct {
	EndpointsApi kube.EndpointsIntf
	MappingApi   kube.MappingIntf
	ServicesMap  map[string]uint16
	EndpointsMap map[string]uint16
	Groups       map[uint16][]string
}

func (sc *ServicesController) Services() map[string]uint16 {
	return sc.ServicesMap
}
func (sc *ServicesController) Endpoints() map[string]uint16 {
	return sc.EndpointsMap
}

func (sc *ServicesController) SyncEndpoints() map[string]string {
	var changed = false

	services, mappingError := sc.MappingApi.UpdateMapping()
	if mappingError == nil {
		if !reflect.DeepEqual(services, sc.Services) {
			changed = true
			sc.ServicesMap = services
		}
	}
	endpoints, endpointsError := sc.EndpointsApi.UpdateEndpoints()
	if endpointsError == nil {
		covertedEndpoints := sc.endpointsMap(endpoints)
		if !reflect.DeepEqual(covertedEndpoints, sc.Endpoints) {
			changed = true
			sc.EndpointsMap = covertedEndpoints
		}
	}

	for _, serviceName := range endpoints {
		if _, ok := services[serviceName]; !ok {
			return endpoints
		}
	}
	if changed {
		sc.Groups = sc.generateGroups()
	}

	return nil
}

func (sc *ServicesController) RandomEndpointByType(serviceType uint16) (string, error) {
	services := sc.Groups[serviceType]
	size := len(services)
	if size > 0 {
		return services[rand.Intn(size)], nil
	} else {
		return "", errors.New("not found")
	}
}

func (sc *ServicesController) generateGroups() map[uint16][]string {
	groups := make(map[uint16][]string)
	for endpoint, _ := range sc.EndpointsMap {
		typeId := sc.EndpointsMap[endpoint]
		if groups[typeId] == nil {
			groups[typeId] = []string{}
		}
		group := groups[typeId]
		if !slices.Contains(group, endpoint) {
			groups[typeId] = append(group, endpoint)
		}
	}
	return groups
}

func (sc *ServicesController) endpointsMap(endpointsKeys map[string]string) map[string]uint16 {
	res := make(map[string]uint16)
	for key, serviceName := range endpointsKeys {
		res[key] = sc.ServicesMap[serviceName]
	}
	return res
}
