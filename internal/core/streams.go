package core

type ValidateResult uint8

const StreamError uint8 = 13

const (
	StreamCollision ValidateResult = iota
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
	wsEndpoint  string
	zmqEndpoint string
}

type RemoveParam struct {
	id      uint32
	message string
}

func NewRouting(sc *ServicesController) *Routing {
	return &Routing{
		streams: make(map[uint32]*StreamConfig),
		remove:  make(chan *RemoveParam, 256),
		sc:      sc,
	}
}

type Routing struct {
	streams map[uint32]*StreamConfig
	remove  chan *RemoveParam
	sc      *ServicesController
}

func (r Routing) removeStreamsByZmq(zmeEndpoint string, message string) { //todo испраивить message на код
	for id, streamConfig := range r.streams {
		if streamConfig.zmqEndpoint == zmeEndpoint {
			r.remove <- &RemoveParam{id, message}
		}
	}
}

func (r Routing) removeStreamsByWs(wsEndpoint string, message string) { //todo испраивить message на код
	for id, streamConfig := range r.streams {
		if streamConfig.wsEndpoint == wsEndpoint {
			r.remove <- &RemoveParam{id, message}
		}
	}
}

func (r Routing) wsRouteNew(messageWrapper *MessagePackage) (*StreamConfig, ValidateResult) {
	var id = messageWrapper.Stream()

	if _, has := r.streams[id]; has {
		return nil, StreamCollision
	} else {
		serviceType := messageWrapper.Service()
		randomZmqConnection, err := r.sc.getRandomServiceConnectionByType(serviceType)
		if err == nil {
			config := &StreamConfig{messageWrapper.ConnectionKey, randomZmqConnection}
			r.streams[id] = config
			return config, Ok
			//	r.sendMessage(messageWrapper, randomZmqConnection)
		} else {
			return nil, EndpointNotFound
		}
	}
}

func (r Routing) wsRoute(messageWrapper *MessagePackage) (*StreamConfig, ValidateResult) {
	var id = messageWrapper.Stream()

	if conf, has := r.streams[id]; has {
		if conf.wsEndpoint != messageWrapper.ConnectionKey {
			return nil, StreamCollision
		} else {
			return conf, Ok
		}
	} else {
		return nil, StreamNotFound
	}

}

func (r Routing) zqmRoute(messageWrapper *MessagePackage) (*StreamConfig, ValidateResult) {
	if conf, has := r.streams[messageWrapper.Stream()]; has {
		if conf.wsEndpoint != messageWrapper.ConnectionKey {
			return nil, AccessDenied
		} else {
			return conf, Ok
		}
	} else {
		return nil, StreamNotFound
	}

}
