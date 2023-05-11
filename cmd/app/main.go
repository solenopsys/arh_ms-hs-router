package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	bl_kubernetes_tools "github.com/solenopsys/bl-kubernetes-tools"
	"math/rand"
	"os"
	"solenopsys.org/zmq_router/internal/conf"
	"solenopsys.org/zmq_router/internal/core"
	"solenopsys.org/zmq_router/internal/io/kube/dev"

	"sync"
	"time"
)

var Mode string

const DEV_MODE = "dev"

func init() {
	flag.StringVar(&Mode, "mode", "", "a string var")
}

func main() {
	flag.Parse()
	fmt.Println("mode:", Mode)

	rand.Seed(time.Now().Unix())

	devMode := Mode == DEV_MODE

	if devMode {
		godotenv.Load("configs/router/local.env")
	}

	host := os.Getenv("server.Host")
	port := os.Getenv("server.Port")
	endpointsPort := os.Getenv("nodes.Port")

	kubeConfig, _ := bl_kubernetes_tools.CreateKubeConfig(devMode)
	var integrator conf.Integrator

	if devMode {
		controller := &core.ServicesController{
			EndpointsApi: dev.NewEndpointsIO(),
			MappingApi:   dev.NewMappingIO(),
			ServicesMap:  make(map[string]uint16),
			EndpointsMap: make(map[string]uint16),
			Groups:       make(map[uint16][]string),
		}

		integrator = conf.Integrator{
			HttpPort:           port,
			HttpHost:           host,
			EndpointPort:       endpointsPort,
			KubeConfig:         kubeConfig,
			ServicesController: controller,
		}
		integrator.SyncLoopInit(nil)
		integrator.InitTest()
	} else {

		integrator = conf.Integrator{
			HttpPort:     port,
			HttpHost:     host,
			EndpointPort: endpointsPort,
			KubeConfig:   kubeConfig,
		}
		integrator.Init()

	}

	integrator.ServicesController.SyncEndpoints()

	var wg sync.WaitGroup //todo перенести это
	wg.Add(1)
	wg.Wait()
}
