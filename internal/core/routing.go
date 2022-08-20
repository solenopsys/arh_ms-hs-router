package core

type ValidateResult uint8

const StreamError uint8 = 13

const (
	StreamCollision ValidateResult = iota
	FirstFlagNotFound
	EndpointNotFound
	StreamNotFound
	AccessDenied
	Ok
)

func (r ValidateResult) String() string {
	return [...]string{
		"StreamCollision",
		"EndpointNotFound",
		"StreamNotFound",
		"AccessDenied",
		"Ok",
	}[r]
}

type StreamConfig struct {
	WsEndpoint  string
	ZmqEndpoint string
}

type RemoveParam struct {
	id      uint32
	message string
}

type Routing struct {
	Routes        map[uint32]*StreamConfig
	Remove        chan *RemoveParam //todo
	RandomService func(serviceId uint16) (string, error)
}

func (r Routing) WsRouteNew(message MessageRoute) (*StreamConfig, ValidateResult) {
	var id = message.Stream()
	if !message.IsFirst() {
		return nil, FirstFlagNotFound
	}
	if _, has := r.Routes[id]; has {
		return nil, StreamCollision
	}
	randomZmqConnection, err := r.RandomService(message.Service())
	if err != nil {
		return nil, EndpointNotFound
	}
	config := &StreamConfig{message.ConnectionKey(), randomZmqConnection}
	r.Routes[id] = config
	return config, Ok
}

func (r Routing) WsRoute(messageWrapper MessageRoute) (*StreamConfig, ValidateResult) {
	var id = messageWrapper.Stream()

	if conf, has := r.Routes[id]; has {
		if conf.WsEndpoint != messageWrapper.ConnectionKey() {
			return nil, StreamCollision
		} else {
			return conf, Ok
		}
	} else {
		return nil, StreamNotFound
	}
}

func (r Routing) ZqmRoute(messageWrapper MessageRoute) (*StreamConfig, ValidateResult) {
	if conf, has := r.Routes[messageWrapper.Stream()]; has {
		if conf.ZmqEndpoint != messageWrapper.ConnectionKey() {
			return nil, AccessDenied
		} else {
			return conf, Ok
		}
	} else {
		return nil, StreamNotFound
	}

}

func (r Routing) SelectRoutesByZmq(endpoint string) []uint32 {
	var res []uint32
	for id, streamConfig := range r.Routes {
		if streamConfig.ZmqEndpoint == endpoint {
			res = append(res, id)
		}
	}
	return res
}

func (r Routing) SelectRoutesByWs(endpoint string) []uint32 {
	var res []uint32
	for id, streamConfig := range r.Routes {
		if streamConfig.WsEndpoint == endpoint {
			res = append(res, id)
		}
	}
	return res
}

func (r Routing) DeleteRoute(stream uint32) {
	delete(r.Routes, stream)
}
