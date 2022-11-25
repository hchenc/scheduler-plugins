package seckill

import (
	"context"
	"fmt"
	gochache "github.com/patrickmn/go-cache"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/scheduler-plugins/pkg/apis/config"
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
	handler  framework.Handle
	podCache *gochache.Cache
}

// order: redis master -> redis slave -> seckill-app
func (s *Seckill) Permit(ctx context.Context, state *framework.CycleState, p *corev1.Pod, nodeName string) (*framework.Status, time.Duration) {
	//master first
	if strings.Contains(p.Name, "seckill-master") && p.Labels["role"] == "master" {
		return nil, 0
	}
	//slave second
	if strings.Contains(p.Name, "seckill-slave") && p.Labels["role"] == "slave" {
		_, ok := s.podCache.Get(p.Name)
		if !ok {
			s.podCache.Set(p.Name, "pod", 0)
			return framework.NewStatus(framework.Wait, "should wait"), 15 * time.Second
		}
		return nil, 0
	}

	if strings.Contains(p.Name, "seckill-mysql") {
		_, ok := s.podCache.Get(p.Name)
		if !ok {
			s.podCache.Set(p.Name, "pod", 0)
			return framework.NewStatus(framework.Wait, "should wait"), 45 * time.Second
		}
		return nil, 0
	}

	if strings.Contains(p.Name, "seckill-statis") {
		_, ok := s.podCache.Get(p.Name)
		if !ok {
			s.podCache.Set(p.Name, "pod", 0)
			return framework.NewStatus(framework.Wait, "should wait"), 45 * time.Second
		}
		return nil, 0
	}

	//if strings.Contains(p.Name, "seckill-app") {
	//	_, ok := s.podCache.Get(p.Name)
	//	if !ok {
	//		s.podCache.Set(p.Name, "pod", 0)
	//		return framework.NewStatus(framework.Wait, "should wait"), 7 * time.Second
	//	}
	//	return nil, 0
	//}

	if strings.Contains(p.Name, "stress-boss") {
		_, ok := s.podCache.Get(p.Name)
		if !ok {
			s.podCache.Set(p.Name, "pod", 0)
			return framework.NewStatus(framework.Wait, "should wait"), 100 * time.Second
		}
		return nil, 0
	}

	if strings.Contains(p.Name, "stress-worker") {
		_, ok := s.podCache.Get(p.Name)
		if !ok {
			s.podCache.Set(p.Name, "pod", 0)
			return framework.NewStatus(framework.Wait, "should wait"), 120 * time.Second
		}
		return nil, 0
	}

	return nil, 0

}

func (s *Seckill) Name() string {
	return Name
}

func New(obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	args, ok := obj.(*config.SeckillArgs)
	if !ok {
		fmt.Println(args)
	}

	fmt.Println("Regis:**********: ")
	s := &Seckill{
		h,
		gochache.New(3*time.Minute, 3*time.Minute),
	}
	return s, nil
}
