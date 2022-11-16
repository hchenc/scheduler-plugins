package seckill

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"strings"
	"time"
)

const Name = "SeckillSort"

var (
	redisSlaveReady = false
	mysqlReady      = false
	statisReady     = false
	appReady        = false
	bossReady       = false
	workerReady     = false
)

//// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//type SeckillSchedulingArgs struct {
//	metav1.TypeMeta
//
//	// KubeConfigPath is the path of kubeconfig.
//	KubeConfigPath string
//}

type Seckill struct {
	handler framework.Handle
}

func (s *Seckill) Reserve(ctx context.Context, state *framework.CycleState, p *corev1.Pod, nodeName string) *framework.Status {
	return nil
}

func (s *Seckill) Unreserve(ctx context.Context, state *framework.CycleState, p *corev1.Pod, nodeName string) {
	s.handler.IterateOverWaitingPods(func(waitingPod framework.WaitingPod) {
		waitingPod.Allow(p.Name)
	})
}

// order: redis master -> redis slave -> seckill-app
func (s *Seckill) Permit(ctx context.Context, state *framework.CycleState, p *corev1.Pod, nodeName string) (*framework.Status, time.Duration) {
	//master first
	if strings.Contains(p.Name, "seckill-master") && p.Labels["role"] == "master" {
		return nil, 0
	}
	//slave second
	if strings.Contains(p.Name, "seckill-slave") && p.Labels["role"] == "slave" {
		return framework.NewStatus(framework.Wait, "should wait"), 30 * time.Second
	}

	if strings.Contains(p.Name, "seckill-mysql") {
		return framework.NewStatus(framework.Wait, "should wait"), 45 * time.Second
	}

	if strings.Contains(p.Name, "seckill-statis") {
		return framework.NewStatus(framework.Wait, "should wait"), 45 * time.Second
	}

	if strings.Contains(p.Name, "seckill-app") {
		return framework.NewStatus(framework.Wait, "should wait"), 60 * time.Second
	}

	if strings.Contains(p.Name, "stress-boss") {
		return framework.NewStatus(framework.Wait, "should wait"), 100 * time.Second
	}

	if strings.Contains(p.Name, "stress-worker") {
		return framework.NewStatus(framework.Wait, "should wait"), 120 * time.Second
	}

	return nil, 0

}

func (s *Seckill) Name() string {
	return Name
}

func New(obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	//args, ok := obj.(SeckillSchedulingArgs)

	fmt.Println("Regis:**********: ")
	return &Seckill{
		h,
	}, nil
}
