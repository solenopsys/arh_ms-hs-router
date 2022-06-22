package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var devMode = os.Getenv("developerMode") == "true"
var host = os.Getenv("server.Host")
var port = os.Getenv("server.Port")
var nodesPort = os.Getenv("nodes.Port")
var address = host + ":" + port
var addr = flag.String("addr", address, "http service address")
var upgrader = websocket.Upgrader{}

func main() {
	println("START")
	rand.Seed(time.Now().Unix())
	wsHub := newWsHub()
	zmqHub := newZqmHub()
	routing := newRouting()
	state := newState(zmqHub, wsHub, routing)

	go state.commandProcessor()
	go state.messageProcessor()
	go zmqHub.updateLoop(state, 1*time.Second)

	http.HandleFunc("/info", getInfo(&state.services))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHub.tryConnectionProcessing(state, w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
