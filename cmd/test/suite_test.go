package test

import (
	"solenopsys.org/zmq_router/internal/conf"
	"solenopsys.org/zmq_router/internal/core"
	"solenopsys.org/zmq_router/internal/io/kube/verify"
	"solenopsys.org/zmq_router/pkg/fortest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntegrationRouter(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "HsRouter Integration Tests")

	BeforeAll(func() {

	})
}

var _ = Describe("Router", func() {
	var integrator *conf.Integrator
	var te *fortest.TestEnv
	BeforeEach(func() {
		te = fortest.NewTestEnv("127.0.0.1:85", "tcp://127.0.0.1")
		servicesController := &core.ServicesController{
			EndpointsApi: verify.NewEndpointsIO(),
			MappingApi:   verify.NewMappingIO(),
			ServicesMap:  make(map[string]uint16),
			EndpointsMap: make(map[string]uint16),
			Groups:       make(map[uint16][]string),
		}

		integrator = &conf.Integrator{
			HttpPort:           "888",
			HttpHost:           "127.0.0.1",
			EndpointPort:       "9999",
			ServicesController: servicesController,
		}

		integrator.InitTest()

		te.CreateEndpoint("dgraph", 8237)
		te.CreateEndpoint("dgraph", 8238)
		te.CreateEndpoint("postgres", 2232)
		te.CreateEndpoint("clickhouse", 3322)

		te.CreateClient(1)
		te.CreateClient(2)
		time.Sleep(1 * time.Second)
	})

	Context("Add routes", func() {
		It("should be empty", func() {
			integrator.ServicesController.SyncEndpoints()
			message := fortest.TestMessage{StreamId: 12312, ServiceId: 221, FunctionId: 2, Body: []byte("bla"), First: true}
			te.SendMessage(message, 1)
			messageSecond := fortest.TestMessage{StreamId: 12312, ServiceId: 221, FunctionId: 2, Body: []byte("bla"), First: false}
			te.SendMessage(messageSecond, 1)
			reseived := <-te.ZmqPool.FromHub
			Expect(message).To(Equal(reseived))

		})
	})

})

//func testConnectService() {
//	model := createModel()
//	model.addService(2001, "service1")
//	sleep()
//	chekService()
//}
//
//func testConnectWs() {
//	model := createModel()
//	model.addClient("client1", "trueKey")
//	sleep()
//	chekClient()
//}
//
//func testAuthErorWs() {
//	model := createModel()
//	model.addClient("client1", "falseKey")
//	sleep()
//	chekClient()
//}
//
//func testMessageWsToService() {
//	model := createModel()
//	model.addService(2001, "service1")
//	model.addClient("client1", "trueKey")
//	const streamId = 212
//	const funcId = 212
//	body := []byte("")
//	model.sendMessage("service1", streamId, funcId, body)
//}
//
//func testMessageWsToServiceAndResponse() {
//	model := createModel()
//	model.addService(2001, "service1")
//	model.addClient("client1", "trueKey")
//	const streamId = 212
//	const funcId = 212
//	body := []byte("")
//	model.sendMessage("service1", streamId, funcId, body)

/*func (te TestEnv) testMessages(n int) {
	for i := 0; i < n; i++ {
		println("NEXT MESSAGE")
		time.Sleep(1 * time.Second)
		body1 := []byte("MessageS1-" + strconv.Itoa(i))
		te.wsPool.ToWs <- &WsMessage{ClientId: 1, Body: te.createMessage(2342, 1, 2, body1, 0 == i)}
		body2 := []byte("MessageS2-" + strconv.Itoa(i))
		te.wsPool.ToWs <- &WsMessage{ClientId: 2, Body: te.createMessage(23422, 0, 3, body2, 0 == i)}
	}
}*/

//}
