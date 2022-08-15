package prod

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"solenopsys.org/zmq_router/internal/io/kube"
)

func NewEndpointsIO(clientset *kubernetes.Clientset, port string) kube.EndpointsIntf {
	e := &EndpointsIO{clientset: clientset}
	_, err := e.UpdateEndpoints()
	if err != nil {
		klog.Error("ERROR GET ENDPOINTS", err)
	}
	return e
}

type EndpointsIO struct {
	clientset *kubernetes.Clientset
	endpoints map[string]string // endpoint : service
	port      string
}

func (e EndpointsIO) Endpoints() map[string]string {
	return e.endpoints
}

func (e EndpointsIO) UpdateEndpoints() (map[string]string, error) {
	endpoints := make(map[string]string)

	pods, err := e.clientset.CoreV1().Pods("default").List(context.Background(), v1.ListOptions{LabelSelector: "type=hStreamNode"})
	if err != nil {
		return nil, err
	}
	for _, pod := range pods.Items {
		ip := pod.Status.PodIP
		endPoint := fmt.Sprintf("tcp://%s:%s", ip, e.port)
		serviceName := pod.Labels["hsServiceName"]
		endpoints[endPoint] = serviceName
	}
	e.endpoints = endpoints

	return endpoints, nil
}
