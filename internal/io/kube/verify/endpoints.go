package verify

import (
	"encoding/json"
	"k8s.io/klog/v2"
	"solenopsys.org/zmq_router/internal/io/kube"
	"solenopsys.org/zmq_router/pkg/utils"
)

func NewEndpointsIO() kube.EndpointsIntf {
	e := &EndpointsIO{}
	_, err := e.UpdateEndpoints()
	if err != nil {
		klog.Error("ERROR GET ENDPOINTS", err)
	}
	return e
}

type EndpointsIO struct {
	endpoints map[string]string // endpoint : service
}

func (e *EndpointsIO) Endpoints() map[string]string {
	return e.endpoints
}

func (e *EndpointsIO) UpdateEndpoints() (map[string]string, error) {
	body, err := utils.ReadBody("/endpoints")
	err = json.Unmarshal([]byte(body), &e.endpoints)
	return e.endpoints, err
}
