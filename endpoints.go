package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
	"strconv"
	"time"
)

func getEndpoints(clientset *kubernetes.Clientset) map[string]string {
	endpoints := make(map[string]string)

	pods, err := clientset.CoreV1().Pods("default").List(context.Background(), v1.ListOptions{LabelSelector: "type=hStreamNode"})
	if err != nil {
		klog.Error("error getting pods: %v\n", err)
	}
	for _, pod := range pods.Items {
		ip := pod.Status.PodIP
		endPoint := fmt.Sprintf("tcp://%s:%s", ip, nodesPort)
		serviceName := pod.Labels["hsServiceName"]
		endpoints[endPoint] = serviceName
	}

	return endpoints
}

func lock(clientset *kubernetes.Clientset, run func(ctx context.Context)) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := uuid.New().String()

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "hstream-lock",
			Namespace: "default",
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				run(ctx)
			},
			OnStoppedLeading: func() {
				klog.Infof("leader lost: %s", id)
			},
			OnNewLeader: func(identity string) {
				klog.Infof("new leader elected: %s", identity)
			},
		},
	})
}

func newEndpoint(endpoints map[string]string, currentMapping map[string]uint16) (map[string]uint16, bool) {
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

func convertMapToInt(stringMap map[string]string) map[string]uint16 {
	numberMap := make(map[string]uint16)
	for key, value := range stringMap {
		parseUint, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			klog.Error("Error parce integer: %v\n", err)
		} else {
			numberMap[key] = uint16(parseUint)
		}
	}
	return numberMap
}

func convertMapToString(numberMap map[string]uint16) map[string]string {
	stringMap := make(map[string]string)
	for key, value := range numberMap {
		stringMap[key] = strconv.Itoa(int(value))
	}
	return stringMap
}

func changeConfigMap(clientSet *kubernetes.Clientset, endpoints map[string]string, s *State) error {

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
		changedData, changed := newEndpoint(endpoints, convertMapToInt(configMapData))
		if changed {
			configMap.Data = convertMapToString(changedData)
			clientSet.CoreV1().ConfigMaps(ns).Create(ctx, &configMap, metav1.CreateOptions{})

			s.incStat("ChangingHsMapping")
		}
	} else {
		data := beforeMap.Data
		intData := convertMapToInt(data)
		changedData, changed := newEndpoint(endpoints, intData)
		if changed {
			configMap.Data = convertMapToString(changedData)
			clientSet.CoreV1().ConfigMaps(ns).Update(ctx, &configMap, metav1.UpdateOptions{})
			s.incStat("ChangingHsMapping")
		}
	}
	return nil
}

func getHsMapping(clientset *kubernetes.Clientset) map[string]uint16 {
	klog.Info("Start...")
	ctx := context.TODO()
	maps, err := clientset.CoreV1().ConfigMaps("default").Get(ctx, "hs-mapping", metav1.GetOptions{})
	if err != nil {
		klog.Error("error getting maps: %v\n", err)
	}

	return convertMapToInt(maps.Data)
}

func configMapUpdateNeeded(clientset *kubernetes.Clientset) bool {
	endpoints := getEndpoints(clientset)
	mapping := getHsMapping(clientset)
	for _, serviceName := range endpoints {
		if _, ok := mapping[serviceName]; !ok {

			return true
		}
	}
	return false
}

func updateConfigMap(clientset *kubernetes.Clientset, s *State) {
	updateNeeded := configMapUpdateNeeded(clientset)
	if updateNeeded {
		run := func(ctx context.Context) {
			endpoints := getEndpoints(clientset)
			err := changeConfigMap(clientset, endpoints, s)
			if err != nil {
				klog.Error("ERROR CHANGE CONFIG...", err)
			}
		}

		lock(clientset, run)
	}

}
