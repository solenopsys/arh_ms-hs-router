package fortest

func NewTestZmqPool() *ZmqServersPool {
	zmq := ZmqServersPool{
		FromHub: make(chan *ZmqMessage, 256),
		ToHub:   make(chan *ZmqMessage, 256),
		Servers: make(map[uint16]*ZmqMockServer),
	}
	go zmq.sendLoop()
	return &zmq
}

func NewTestWsPool() *WsClientsPool {
	pool := WsClientsPool{
		FromWs:            make(chan *WsMessage, 256),
		ToWs:              make(chan *WsMessage, 256),
		ClientConnections: make(map[uint16]*WsMockClient),
	}
	go pool.sendLoop()
	return &pool
}

func NewHttp(address string, defaultServiceUrl string) *MockHttp {
	http := MockHttp{
		address:    address,
		Services:   make(map[string]uint16),
		Endpoints:  make(map[uint16]uint16),
		defaultUrl: defaultServiceUrl,
	}
	go http.start()
	return &http
}

func NewTestEnv(address string, defaultServiceUrl string) *TestEnv {
	return &TestEnv{
		ZmqPool: NewTestZmqPool(),
		WsPool:  NewTestWsPool(),
		Http:    NewHttp(address, defaultServiceUrl),
		zmqHost: defaultServiceUrl,
	}
}
