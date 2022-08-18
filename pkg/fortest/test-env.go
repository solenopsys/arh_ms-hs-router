package fortest

import (
	"encoding/binary"
	"solenopsys.org/zmq_router/internal/core"
)

type TestMessage struct {
	StreamId   uint32
	ServiceId  uint16
	FunctionId uint8
	Body       []byte
	First      bool
}

func (m TestMessage) ToBytes() []byte {

	header := make([]byte, 8)
	binary.BigEndian.PutUint32(header[0:], m.StreamId)
	if m.First {
		header[4] = core.FirstFrame
	} else {
		header[4] = core.OrdinaryFrame
	}

	header[5] = m.FunctionId
	binary.BigEndian.PutUint16(header[6:], m.ServiceId)

	bytes := append(header[:], m.Body[:]...)
	return bytes
}

type TestEnv struct {
	ZmqPool        *ZmqServersPool // clients
	WsPool         *WsClientsPool  // clients
	Http           *MockHttp       // clients
	serviceCounter uint16
	clientsCounter uint16
	zmqHost        string
	wsAddress      string
}

func (te *TestEnv) CreateEndpoint(serviceName string, port uint16) *ZmqMockServer {
	server := te.ZmqPool.CreateServer(te.zmqHost, port)
	if _, ok := te.Http.Services[serviceName]; !ok {
		te.Http.Services[serviceName] = te.serviceCounter
		te.serviceCounter++
	}
	te.Http.Endpoints[port] = te.Http.Services[serviceName]

	return server
}

func (te *TestEnv) CreateClient(clientId uint16) (*WsMockClient, error) {
	return te.WsPool.CreateClient(te.wsAddress, clientId)
}

func (te *TestEnv) SendMessage(m TestMessage, clientId uint16) {
	te.WsPool.ToWs <- &WsMessage{clientId, m.ToBytes()}
}
