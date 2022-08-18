package kube

type EndpointsIntf interface {
	UpdateEndpoints() (map[string]string, error)
}
type MappingIntf interface {
	UpdateMapping() (map[string]uint16, error)
}
