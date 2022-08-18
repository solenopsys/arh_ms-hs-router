package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"solenopsys.org/zmq_router/internal/conf"
	"solenopsys.org/zmq_router/pkg/kube"
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
		godotenv.Load("configs/router/dev.env")
	}

	host := os.Getenv("server.Host")
	port := os.Getenv("server.Port")
	endpointsPort := os.Getenv("nodes.Port")

	kubeConfig := kube.CreateKubeConfig(devMode)

	integrator := conf.Integrator{
		HttpPort:     port,
		HttpHost:     host,
		EndpointPort: endpointsPort,
		KubeConfig:   kubeConfig,
	}
	integrator.Init()
}
