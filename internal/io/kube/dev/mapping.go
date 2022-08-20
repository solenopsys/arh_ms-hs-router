package dev

import (
	kube2 "solenopsys.org/zmq_router/internal/io/kube"
)

func NewMappingIO() kube2.MappingIntf {
	return &MappingIO{}
}

type MappingIO struct {
	mapping map[string]uint16 // endpoint : service
}

func (e MappingIO) Mapping() map[string]uint16 {
	mapping, _ := e.UpdateMapping()
	return mapping
}

func (e MappingIO) UpdateMapping() (map[string]uint16, error) {
	return map[string]uint16{
		//	"alexstorm-hsm-installer": 1,
		"alexstorm-hsm-dgraph": 2,
	}, nil
}
