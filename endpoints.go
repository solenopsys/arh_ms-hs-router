package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"log"
	"strconv"
)

func decodeEndpoints(body []byte) (map[string]uint16, map[string]uint16) {
	var podsResponse PodsResponse

	services := make(map[string]uint16)
	endpoints := make(map[string]uint16)
	err := json.Unmarshal(body, &podsResponse)
	if err != nil {
		//return //todo обработать
	}

	for _, item := range podsResponse.Items {
		ip := item.Status.PodIP
		endPoint := fmt.Sprintf("tcp://%s:%s", ip, nodesPort)
		println(item.Metadata.Name)
		labels := item.Metadata.Labels
		parseInt, err := strconv.ParseInt(labels["hsID"], 10, 16)
		if err != nil {
			//return //todo обработать
		}
		services[labels["hsServiceName"]] = uint16(parseInt)
		endpoints[endPoint] = uint16(parseInt)
	}

	return services, endpoints
}

func getEndpoints() (map[string]uint16, map[string]uint16) {
	uri := "/api/v1/pods?labelSelector=type%3DhStreamNode"

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientFor, err := rest.HTTPClientFor(config)
	if err != nil {
		panic(err.Error())
	}

	url := fmt.Sprintf("%s%s", config.Host, uri)

	println("URL:")
	println(url)
	resp, err := clientFor.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	println("BODY:")
	//	println(string(body))
	return decodeEndpoints(body)
}

func getEndpointsDeveloperMode() (map[string]uint16, map[string]uint16) {
	services := map[string]uint16{
		"clickhouse": 3,
		"dgraph":     2,
		"postgres":   1,
		"kubernetes": 4,
		"git":        5,
	}
	endpoints := map[string]uint16{
		"tcp://localhost:5559": 1,
		"tcp://localhost:5558": 2,
		"tcp://localhost:5556": 3,
		"tcp://localhost:5557": 4,
		"tcp://localhost:5575": 5,
	}
	return services, endpoints
}

type PodsResponse struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name   string            `json:"name"`
			Labels map[string]string `json:"labels"`
		} `json:"metadata"`
		Status struct {
			PodIP string `json:"podIP"`
		} `json:"status"`
	} `json:"items"`
}
