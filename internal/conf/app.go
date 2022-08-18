package conf

import (
	"github.com/go-zeromq/zmq4"
	"k8s.io/client-go/kubernetes"
	"solenopsys.org/zmq_router/internal/core"
	"solenopsys.org/zmq_router/internal/io/http"
	"solenopsys.org/zmq_router/internal/io/kube/prod"
	"solenopsys.org/zmq_router/internal/io/ws"
	"solenopsys.org/zmq_router/internal/io/zmq"
	"solenopsys.org/zmq_router/pkg/utils"
	"time"
)

type Integrator struct {
	statistic          *utils.Statistic
	KubeConfig         *kubernetes.Clientset
	ServicesController *core.ServicesController
	routing            *core.Routing
	bridge             *core.MessagesBridge
	wsHub              *ws.Hub
	zmqHub             *zmq.Hub
	httpServer         *http.HttpController
	HttpHost           string
	HttpPort           string
	EndpointPort       string
}

func (c *Integrator) initStatistic() {
	c.statistic = &utils.Statistic{
		Values:    make(map[string]uint16),
		Increment: make(chan string, 256),
	}
	go c.statistic.UpdateLoop()
}

func (c *Integrator) initServiceController() {
	c.ServicesController = &core.ServicesController{
		EndpointsApi: prod.NewEndpointsIO(c.KubeConfig, c.EndpointPort),
		MappingApi:   prod.NewMappingIO(c.KubeConfig),
		ServicesMap:  make(map[string]uint16),
		EndpointsMap: make(map[string]uint16),
		Groups:       make(map[uint16][]string),
	}

	go func() {
		for range time.Tick(10 * time.Second) {
			c.ServicesController.SyncEndpoints()
		}
	}()
}

func (c *Integrator) initRouting() {
	c.routing = &core.Routing{
		Routes: make(map[uint32]*core.StreamConfig),
		Remove: make(chan *core.RemoveParam, 256),
		RandomService: func(serviceId uint16) (string, error) {
			return c.ServicesController.RandomEndpointByType(serviceId)
		},
	}
}

func (c *Integrator) initWsHub() {
	c.wsHub = &ws.Hub{
		Connections: make(map[string]*ws.Connection),
		Commands:    make(chan *ws.Command, 256),
		Events:      make(chan *ws.Event, 256),
		Input:       make(chan *ws.Message, 256),
		Output:      make(chan *ws.Message, 256),
	}
	go c.wsHub.SendToWsMessageProcessing()
}

func (c *Integrator) initZmqHub() {
	c.zmqHub = &zmq.Hub{
		Connections: make(map[string]*zmq4.Socket),
		Commands:    make(chan *zmq.Command, 256),
		Events:      make(chan *zmq.Event, 256),
		Input:       make(chan *zmq.Message, 256),
		Output:      make(chan *zmq.Message, 256),
	}
	go c.zmqHub.CommandProcessing()

}

func (c *Integrator) initMessagesBridge() {
	c.bridge = &core.MessagesBridge{
		WsHub:     c.wsHub,
		ZmqHub:    c.zmqHub,
		Router:    c.routing,
		Statistic: c.statistic,
	}
}

func (c *Integrator) initHttp() {
	c.httpServer = &http.HttpController{
		Services:  c.ServicesController.Services,
		Statistic: c.statistic.GetValues,
		Endpoints: c.ServicesController.Endpoints,
		WsHub:     c.wsHub,
	}

	go c.httpServer.StartServer("/hs/", c.HttpHost, c.HttpPort)
}

func (c *Integrator) Init() {
	c.initStatistic()
	c.initServiceController()
	c.initRouting()
	c.initWsHub()
	c.initZmqHub()
	c.initMessagesBridge()
	c.initHttp()
}

func (c *Integrator) InitTest() {
	c.initStatistic()
	c.initRouting()
	c.initWsHub()
	c.initZmqHub()
	c.initMessagesBridge()
	c.initHttp()
}

func (c *Integrator) start() {

	go c.wsEventProcessing()
	go c.zmqEventProcessing()
	//	go cd.SyncConnectionScheduler()

}

func (c *Integrator) wsEventProcessing() {
	for {
		event := <-c.wsHub.GetEvents()
		if event.EventType == ws.OnConnected {
			c.statistic.Increment <- "WsConnect"
		}
	}
}

func (c *Integrator) zmqEventProcessing() {
	for {
		event := <-c.zmqHub.GetEvents()
		if event.EventType == zmq.OnConnected {
			c.statistic.Increment <- "ZmqConnect"
		}
	}
}
