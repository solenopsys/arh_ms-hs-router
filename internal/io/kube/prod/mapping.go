package prod

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kube2 "solenopsys.org/zmq_router/internal/io/kube"
	"solenopsys.org/zmq_router/pkg/kube"
	"solenopsys.org/zmq_router/pkg/utils"
)

type MappingIO struct {
	clientset   *kubernetes.Clientset
	endpointsIO *EndpointsIO
	kubeLock    *kube.KubeLock
	mapping     map[string]uint16
}

func NewMappingIO(clientset *kubernetes.Clientset) kube2.MappingIntf {
	lock := &kube.KubeLock{clientset, 1, 5, 10}
	e := &MappingIO{clientset: clientset, kubeLock: lock}
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

func (k MappingIO) configMapUpdateNeeded() (bool, error) {
	endpoints, err := k.endpointsIO.UpdateEndpoints()
	if err != nil {
		return false, err
	}
	mapping, err := k.UpdateMapping()
	if err != nil {
		return false, err
	}
	for _, serviceName := range endpoints {
		if _, ok := mapping[serviceName]; !ok {
			return true, nil
		}
	}
	return false, nil
}

func (k MappingIO) Mapping() map[string]uint16 {
	return k.mapping
}

func (k MappingIO) UpdateConfigMap() {

	updateNeeded, err := k.configMapUpdateNeeded()
	if err == nil {

		if updateNeeded {
			k.kubeLock.Lock("hstream-lock", func(ctx context.Context) {
				err = k.changeConfigMap(k.endpointsIO.Endpoints())
				if err != nil {
					klog.Error("ERROR CHANGE CONFIG ", err)
				}
			})
		}
	}

}

func (k MappingIO) changeConfigMap(endpoints map[string]string) error {
	clientSet := k.clientset

	configMapData := make(map[string]string, 0)

	ctx := context.TODO()

	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hs-mapping",
			Namespace: "default",
		},
		Data: configMapData,
	}

	ns := "default"
	if beforeMap, err := clientSet.CoreV1().ConfigMaps(ns).Get(ctx, "hs-mapping", metav1.GetOptions{}); errors.IsNotFound(err) {
		changedData, changed := k.newEndpoint(endpoints, utils.ConvertMapToInt(configMapData))
		if changed {
			configMap.Data = utils.ConvertMapToString(changedData)
			clientSet.CoreV1().ConfigMaps(ns).Create(ctx, &configMap, metav1.CreateOptions{})
		}
	} else {
		data := beforeMap.Data
		intData := utils.ConvertMapToInt(data)
		changedData, changed := k.newEndpoint(endpoints, intData)
		if changed {
			configMap.Data = utils.ConvertMapToString(changedData)
			clientSet.CoreV1().ConfigMaps(ns).Update(ctx, &configMap, metav1.UpdateOptions{})
		}
	}
	return nil
}

func (k MappingIO) newEndpoint(endpoints map[string]string, currentMapping map[string]uint16) (map[string]uint16, bool) {
	newMapping := make(map[string]uint16)

	var max uint16 = 0
	var changed = false

	for _, serviceId := range currentMapping {
		if serviceId > max {
			max = serviceId
		}
	}

	for _, serviceName := range endpoints {
		if val, ok := currentMapping[serviceName]; ok {
			newMapping[serviceName] = val
		} else {
			max++
			newMapping[serviceName] = max
			changed = true
		}
	}
	return newMapping, changed
}
