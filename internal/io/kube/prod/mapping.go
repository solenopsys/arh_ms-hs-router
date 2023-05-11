package prod

import (
	"context"
	"github.com/solenopsys/bl-kubernetes-tools"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"solenopsys.org/zmq_router/pkg/utils"
)

type MappingIO struct {
	clientset *kubernetes.Clientset
	kubeLock  *bl_kubernetes_tools.KubeLock
	mapping   map[string]uint16
}

func NewMappingIO(clientset *kubernetes.Clientset) MappingIO {
	lock := &bl_kubernetes_tools.KubeLock{clientset, 1, 5, 10}
	e := MappingIO{clientset: clientset, kubeLock: lock}
	_, err := e.UpdateMapping()
	if err != nil {
		klog.Error("ERROR GET HS MAPPING", err)
	}
	return e
}

func (k MappingIO) UpdateMapping() (map[string]uint16, error) {
	klog.Info("Start...")
	ctx := context.TODO()
	maps, err := k.clientset.CoreV1().ConfigMaps("default").Get(ctx, "hs-mapping", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	k.mapping = utils.ConvertMapToInt(maps.Data)

	return k.mapping, nil
}

func (k MappingIO) Mapping() map[string]uint16 {
	return k.mapping
}
