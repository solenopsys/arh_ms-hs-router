package core

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"solenopsys.org/zmq_router/internal/io/http"
	"solenopsys.org/zmq_router/internal/io/kube/verify"
	"solenopsys.org/zmq_router/internal/io/ws"
	"solenopsys.org/zmq_router/internal/io/zmq"
	"solenopsys.org/zmq_router/pkg/utils"
	"time"
)

var Host string
var Mode string
var Port string
var NodesPort string

func init() {
	// example with short version for long flag
	flag.StringVar(&Mode, "mode", "", "a string var")
}

const DEV_MODE = "dev"
const VERIFY_MODE = "verify"
const PROD_MODE = "prod"

func NewServicesController(zmqHub zmq.IO, endpointPort string) *ServicesController {

	//	if Mode == DEV_MODE {
	//kubeConfig := k.CreateKubeConfig(true)
	//return &ServicesController{
	//	services:     make(map[string]uint16),
	//	endpoints:    make(map[string]uint16),
	//	endpointsApi: prod.NewEndpointsIO(kubeConfig, endpointPort),
	//	mappingApi:   prod.NewMappingIO(kubeConfig),
	//	zmqHub:       zmqHub,
	//}
	//}
	//	if Mode == VERIFY_MODE {
	return &ServicesController{
		conf: &ServicesConf{
			services:  make(map[string]uint16),
			endpoints: make(map[string]uint16),
			groups:    make(map[uint16][]string),
		},
		endpointsApi: verify.NewEndpointsIO(),
		mappingApi:   verify.NewMappingIO(),
		zmqHub:       zmqHub,
	}
	//	}

	//kubeConfig := k.CreateKubeConfig(false)
	//return &ServicesController{
	//	services:     make(map[string]uint16),
	//	endpoints:    make(map[string]uint16),
	//	endpointsApi: prod.NewEndpointsIO(kubeConfig, endpointPort),
	//	mappingApi:   prod.NewMappingIO(kubeConfig),
	//	zmqHub:       zmqHub,
	//}
}

func Run() {
	flag.Parse()
	fmt.Println("mode:", Mode)
	if Mode != PROD_MODE {
		godotenv.Load("configs/router/" + Mode + ".env")
	}

	Host = os.Getenv("server.Host")
	Port = os.Getenv("server.Port")
	NodesPort = os.Getenv("nodes.Port")

	rand.Seed(time.Now().Unix())

	zmqHub := zmq.NewHub()

	hub := ws.NewHub()

	statistic := utils.NewStatistic()
	servicesController := NewServicesController(zmqHub, NodesPort)
	routing := NewRouting(servicesController)
	bridge := &MessagesBridge{wsHub: hub, zmqHub: zmqHub, router: routing, statistic: statistic}

	coordinator := &Coordinator{
		statistic: statistic,
		wsHub:     hub,
		zmqHub:    zmqHub,
	}

	go servicesController.UpdateEndpointsLoop(10 * time.Second)
	coordinator.start()

	statValues := statistic.Values()
	println("BLA")
	server := http.HttpController{
		Services: func() map[string]uint16 {
			return servicesController.conf.services
		},
		Statistic: &statValues,
		Endpoints: nil,
		WsHub:     hub,
	}

	bridge.init()
	server.StartServer("/hs/", Host, Port)
	//
	//go state.statProcessor()
	//go state.commandProcessor()
	//go state.messageProcessor()
	//go zmqHub.UpdateEndpointsLoop(state, 10*time.Second)
	//go kube.UpdateConfigMapLoop(state, 10*time.Second)

	// startServer
}

type Coordinator struct {
	wsHub     ws.IO
	zmqHub    zmq.IO
	statistic *utils.Statistic
}

// загрузка списка сервисов
// обновление конфигмапов

func (cd Coordinator) start() {

	go cd.wsEventProcessing()
	go cd.zmqEventProcessing()
	//	go cd.SyncConnectionScheduler()

}

func (cd Coordinator) wsEventProcessing() {
	for {
		event := <-cd.wsHub.Events()
		if event.EventType == ws.OnConnected {
			cd.statistic.Increment <- "WsConnect"
		}
	}
}

func (cd Coordinator) zmqEventProcessing() {
	for {
		event := <-cd.zmqHub.Events()
		if event.EventType == zmq.OnConnected {
			cd.statistic.Increment <- "ZmqConnect"
		}
	}
}
