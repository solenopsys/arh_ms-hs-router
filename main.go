package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var devMode = os.Getenv("developerMode") == "true"

var host = os.Getenv("server.Host")
var port = os.Getenv("server.Port")
var nodesPort = os.Getenv("nodes.Port")
var address = host + ":" + port
var addr = flag.String("addr", address, "http service address")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //todo domains politic
	},
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getCubeConfig(devMode bool) (*rest.Config, error) {
	if devMode {
		var kubeconfigFile = os.Getenv("kubeconfigPath")
		kubeConfigPath := filepath.Join(kubeconfigFile)
		klog.Infof("Using kubeconfig: %s\n", kubeConfigPath)

		kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			klog.Error("error getting Kubernetes config: %v\n", err)
			os.Exit(1)
		}

		return kubeConfig, nil
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		return config, nil
	}
}

// ok
func updateConfigMapLoop(state *State, duration time.Duration) {
	for range time.Tick(duration) {
		id := "UpdateConfigMap"
		state.command <- &Command{key: id, messageType: UpdateConfigMap}
	}
}

func main() {
	klog.Infof("START")
	klog.Infof("DEV MODE ", devMode)

	forConfig := createKubeConfig()

	rand.Seed(time.Now().Unix())
	wsHub := newWsHub()
	zmqHub := newZqmHub()
	routing := newRouting()
	state := newState(zmqHub, wsHub, routing, forConfig)

	go state.statProcessor()
	go state.commandProcessor()
	go state.messageProcessor()
	go zmqHub.updateEndpointsLoop(state, 10*time.Second)
	go updateConfigMapLoop(state, 120*time.Second)

	http.HandleFunc("/info", getInfo(&state.services))
	http.HandleFunc("/stat", getStat(state, routing))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		wsHub.tryConnectionProcessing(state, w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func createKubeConfig() *kubernetes.Clientset {
	config, err := getCubeConfig(devMode)
	if err != nil {
		klog.Info("Config init error...", err)
		os.Exit(1)
	}
	forConfig, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Info("Config init error...", err)
		os.Exit(1)
	}
	return forConfig
}
