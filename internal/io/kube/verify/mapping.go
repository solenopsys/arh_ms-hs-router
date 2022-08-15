package verify

import (
	"encoding/json"
	"k8s.io/klog/v2"
	kube2 "solenopsys.org/zmq_router/internal/io/kube"
	"solenopsys.org/zmq_router/pkg/utils"
)

func NewMappingIO() kube2.MappingIntf {
	m := &MappingIO{}
	_, err := m.UpdateMapping()
	if err != nil {
		klog.Error("ERROR GET MAPPING", err)
	}
	return m
}

type MappingIO struct {
	mapping map[string]uint16 // endpoint : service
}

func (e *MappingIO) Mapping() map[string]uint16 {
	return e.mapping
}

func (e *MappingIO) UpdateMapping() (map[string]uint16, error) {
	body, err := utils.ReadBody("/mapping")
	err = json.Unmarshal([]byte(body), &e.mapping)
	return e.mapping, err
}
