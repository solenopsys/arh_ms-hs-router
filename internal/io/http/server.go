package http

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"solenopsys.org/zmq_router/internal/io/ws"
	"solenopsys.org/zmq_router/pkg/utils"
)

type HttpController struct {
	WsHub     *ws.Hub
	Services  func() map[string]uint16
	Statistic func() map[string]uint16
	Endpoints func() map[string]uint16
}

func (c HttpController) StartServer(prefix string, host string, port string) {
	http.HandleFunc(prefix+"services", c.ServicesHandler())
	http.HandleFunc(prefix+"statistic", c.StatisticHandler())
	http.HandleFunc(prefix+"websocket", c.ConnectionFunction(func(token string) (uint16, error) {
		return 1, nil
	}))

	addr := host + ":" + port
	klog.Infof("Listen addrss %", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (c HttpController) ServicesHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.SetCorsHeaders(w)
		marshal, err := json.Marshal(c.Services())
		if err != nil {
			klog.Error("Get info error:", err)
		} else {
			w.Write(marshal)
		}
	}
}

func (c HttpController) StatisticHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var statMap = map[string]map[string]uint16{}

		statMap["statistic"] = c.Statistic()
		statMap["services"] = c.Services()
		statMap["endpoints"] = c.Endpoints()
		stream := map[string]uint16{}
		//	stream["stream"] = uint16(len(*c.streams))
		statMap["streams"] = stream

		klog.Infof("Stat request:", statMap)
		utils.SetCorsHeaders(w)
		marshal, err := json.Marshal(&statMap)
		if err != nil {
			klog.Error("Get info error:", err)
		} else {
			klog.Infof("Stat request in json:", string(marshal))

			w.Write(marshal)
		}
	}
}

func (c HttpController) ConnectionFunction(getAuth func(token string) (uint16, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// r.Header.Get("token")
		token := r.URL.Query().Get("token")
		if token != "" {
			userId, err := getAuth(token)
			if err == nil {
				utils.SetCorsHeaders(w)
				connection, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					klog.Error("Wpgrade ws:", err)
				} else {
					c.WsHub.AddConnection(connection, userId)
				}

			} else {
				w.WriteHeader(403)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //todo domains politic
	},
}
