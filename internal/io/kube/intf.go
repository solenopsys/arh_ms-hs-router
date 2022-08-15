package kube

type EndpointsIntf interface {
	UpdateEndpoints() (map[string]string, error)
	Endpoints() map[string]string
}
type MappingIntf interface {
	UpdateMapping() (map[string]uint16, error)
	Mapping() map[string]uint16
}
