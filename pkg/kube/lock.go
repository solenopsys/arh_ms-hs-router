package kube

import (
	"context"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
	"time"
)

type KubeLock struct {
	Clientset     *kubernetes.Clientset
	RetryPeriod   time.Duration
	RenewDeadline time.Duration
	LeaseDuration time.Duration
}

func (k KubeLock) Lock(lockName string, run func(ctx context.Context)) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := uuid.New().String()

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockName,
			Namespace: "default",
		},
		Client: k.Clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   k.LeaseDuration * time.Second,
		RenewDeadline:   k.RenewDeadline * time.Second,
		RetryPeriod:     k.RetryPeriod * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				klog.Infof("leader lost: %s", id)
			},
			OnNewLeader: func(identity string) {
				klog.Infof("new leader elected: %s", identity)
			},
		},
	})
}
