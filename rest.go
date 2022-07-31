package main

import (
	"encoding/json"
	"k8s.io/klog/v2"
	"net/http"
)

func getInfo(sm *map[string]uint16) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w)
		marshal, err := json.Marshal(sm)
		if err != nil {
			klog.Error("Get info error:", err)
		} else {
			w.Write(marshal)
		}
	}
}

func getStat(sm *State, routing *RoutingMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var statMap = map[string]map[string]uint16{}

		stream := map[string]uint16{}
		statMap["statistic"] = sm.statistic
		statMap["services"] = sm.services
		statMap["endpoints"] = sm.endpoints
		stream["Ws"] = uint16(len(routing.wsConnections))
		stream["Zmq"] = uint16(len(routing.zmqConnections))
		statMap["streams"] = stream

		klog.Infof("Stat request:", statMap)
		setCorsHeaders(w)
		marshal, err := json.Marshal(&statMap)
		if err != nil {
			klog.Error("Get info error:", err)
		} else {
			klog.Infof("Stat request in json:", string(marshal))

			w.Write(marshal)
		}
	}
}
