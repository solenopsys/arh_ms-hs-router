package fortest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type MockHttp struct {
	address    string
	Services   map[string]uint16
	Endpoints  map[uint16]uint16
	defaultUrl string
}

func (m MockHttp) start() {
	http.HandleFunc("/endpoints", m.endpointsFunction())
	http.HandleFunc("/mapping", m.mappingFunction())
	println("LISTEN ", m.address)
	log.Fatal(http.ListenAndServe(m.address, nil))
}

func (m MockHttp) endpointsFunction() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := make(map[string]string)
		for port, serviceId := range m.Endpoints {
			portStr := strconv.FormatInt(int64(port), 10)
			for serviceName, id := range m.Services {
				if serviceId == id {
					res[m.defaultUrl+":"+portStr] = serviceName
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

func (m MockHttp) mappingFunction() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		marshal, err := json.Marshal(m.Services)
		if err != nil {
			log.Panic("Get info error:", err)
		} else {
			w.Write(marshal)
		}
	}
}
