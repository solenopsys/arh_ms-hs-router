package dev

import (
	"k8s.io/klog/v2"
	"solenopsys.org/zmq_router/internal/io/kube"
)

func NewEndpointsIO() kube.EndpointsIntf {
	e := EndpointsIO{}
	_, err := e.UpdateEndpoints()
	if err != nil {
		klog.Error("ERROR GET ENDPOINTS", err)
	}
	return &e
}

type EndpointsIO struct {
}

func (e EndpointsIO) Endpoints() map[string]string {
	endpoints, _ := e.UpdateEndpoints()
	return endpoints
}

func (e EndpointsIO) UpdateEndpoints() (map[string]string, error) {
	return map[string]string{
		//"tcp://192.168.122.29:5561": "alexstorm-hsm-installer",
		"tcp://192.168.122.29:9999": "alexstorm-hsm-dgraph",
	}, nil
}
