package main

import (
	"encoding/binary"
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"solenopsys.org/zmq_router/internal/core"
	"solenopsys.org/zmq_router/tests/integrity/testio"
	"strconv"
	"time"
)

func endpoints() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := make(map[string]string)
		for port, serviceId := range model.endpoints {
			portStr := strconv.FormatInt(int64(port), 10)
			for serviceName, id := range model.services {
				if serviceId == id {
					res[DefaultServiceUrl+":"+portStr] = serviceName
				}

			}
		}
		marshal, err := json.Marshal(res)
		if err != nil {
			log.Panic("Get info error:", err)
		} else {
			w.Write(marshal)
		}
	}
}

func mapping() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		marshal, err := json.Marshal(model.services)
		if err != nil {
			log.Panic("Get info error:", err)
		} else {
			w.Write(marshal)
		}
	}
}

type Model struct {
	services  map[string]uint16      // service name : service id
	endpoints map[uint16]uint16      // port : service id
	zmqPool   *testio.ZmqServersPool // clients
	wsPool    *testio.WsClientsPool  // clients
}

func createEndpoint(serviceName string, port uint16) *testio.ZmqServer {
	server := model.zmqPool.CreateServer(DefaultServiceUrl, port)

	if _, ok := model.services[serviceName]; !ok {
		model.services[serviceName] = serviceCounter
		serviceCounter++
	}
	model.endpoints[port] = model.services[serviceName]

	return server
}

func createClient(clientId uint16) (*testio.WsClient, error) {
	//currentId := clientsCounter
	return model.wsPool.CreateClient(DefaultRouterWsUrl, clientId)
	//	serviceCounter++

	//	 return currentId
}

func init() {
	// example with short version for long flag
	godotenv.Load("configs/test/.env")

	DefaultServiceUrl = os.Getenv("zmqServerHost")
	DefaultRouterWsUrl = os.Getenv("wsAddress")
}

var DefaultServiceUrl string
var DefaultRouterWsUrl string
var model *Model
var serviceCounter uint16 = 0
var clientsCounter uint16 = 0

func main() {

	model = &Model{
		services:  make(map[string]uint16),
		endpoints: make(map[uint16]uint16),
		zmqPool:   testio.NewTestZmqPool(),
		wsPool:    testio.NewTestWsPool(),
	}

	http.HandleFunc("/endpoints", endpoints())
	http.HandleFunc("/mapping", mapping())

	createEndpoint("dgraph", 8237)
	createEndpoint("dgraph", 8238)
	createEndpoint("postgres", 2232)
	createEndpoint("clickhouse", 3322)

	createClient(1)
	createClient(2)

	go testMessages()
	go printZmqMessages()

	addr := "127.0.0.1:85"
	println("LISTEN ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func createMessage(streamId uint32, serviceId uint16, functionId uint8, body []byte, first bool) []byte {

	header := make([]byte, 8)
	binary.BigEndian.PutUint32(header[0:], streamId)
	if first {
		header[4] = core.FirstFrame
	} else {
		header[4] = core.OrdinaryFrame
	}

	header[5] = functionId
	binary.BigEndian.PutUint16(header[6:], serviceId)

	return append(header[:], body[:]...)
}

func testMessages() {
	for i := 0; i < 100; i++ {
		println("NEXT MESSAGE")
		time.Sleep(1 * time.Second)
		body1 := []byte("MessageS1-" + strconv.Itoa(i))
		model.wsPool.ToWs <- &testio.WsMessage{ClientId: 1, Body: createMessage(2342, 1, 2, body1, 0 == i)}
		body2 := []byte("MessageS2-" + strconv.Itoa(i))
		model.wsPool.ToWs <- &testio.WsMessage{ClientId: 2, Body: createMessage(23422, 0, 3, body2, 0 == i)}
	}
}

func printZmqMessages() {
	for {
		message := <-model.zmqPool.FromHub

		println("FROM HUB-- -------", string(message.Message))
		println("PORT", strconv.FormatUint(uint64(message.Port), 10))
	}

}

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
//}
