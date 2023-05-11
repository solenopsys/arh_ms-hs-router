package prod

import (
	"context"
	"github.com/solenopsys/bl-kubernetes-tools"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"solenopsys.org/zmq_router/pkg/utils"
)

func NewChangerConfigmap(clientset *kubernetes.Clientset) ChangerConfigmapIO {
	lock := &bl_kubernetes_tools.KubeLock{clientset, 1, 3, 6}
	return ChangerConfigmapIO{clientset: clientset, kubeLock: lock}
}

type ChangerConfigmapIO struct {
	clientset *kubernetes.Clientset
	kubeLock  *bl_kubernetes_tools.KubeLock
}

func (k ChangerConfigmapIO) UpdateConfigMap(endpoints map[string]string) {
	klog.Infof("START UPDATE CONFIGMAP BY ENDPOINTS  ", endpoints)
	k.kubeLock.Lock("hstream-lock", func(ctx context.Context) {
		err := k.changeConfigMap(endpoints)
		if err != nil {
			klog.Error("ERROR CHANGE CONFIG ", err)
		}
	})
}

func (k ChangerConfigmapIO) changeConfigMap(endpoints map[string]string) error {
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

func (k ChangerConfigmapIO) newEndpoint(endpoints map[string]string, currentMapping map[string]uint16) (map[string]uint16, bool) {
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
