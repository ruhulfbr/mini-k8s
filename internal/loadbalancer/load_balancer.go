package loadbalancer

import (
	"net/http"
	"sync"

	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type LoadBalancer struct {
	podRepo *repositories.PodRepository
	counter uint64
	mu      sync.Mutex
	servers map[string]*http.Server
}

func NewLoadBalancer(podRepo *repositories.PodRepository) *LoadBalancer {
	return &LoadBalancer{
		podRepo: podRepo,
		servers: make(map[string]*http.Server),
	}
}

//
///*
//Start bootstraps load balancers for all existing services.
//This is useful on controller restart.
//*/
//func (lb *LoadBalancer) Start() {
//	services, err := lb.podRepo.ListPodsGroupedByService(context.Background())
//	if err != nil || len(services) == 0 {
//		log.Println("[LB] no services found to start")
//		return
//	}
//
//	for _, svc := range services {
//		go lb.StartServiceListener(svc.Service, svc.Port)
//	}
//}
//
///*
//StartServiceListener starts (or skips if already running)
//a load balancer for a single service.
//*/
//func (lb *LoadBalancer) StartServiceListener(service string, port int) {
//	lb.mu.Lock()
//	if _, exists := lb.servers[service]; exists {
//		lb.mu.Unlock()
//		log.Printf("[LB] already running for service=%s", service)
//		return
//	}
//
//	server := &http.Server{
//		Addr:    fmt.Sprintf(":%d", port),
//		Handler: lb.serviceHandler(service),
//	}
//
//	lb.servers[service] = server
//	lb.mu.Unlock()
//
//	log.Printf("[LB] started for service=%s on port=%d", service, port)
//
//	go func() {
//		if err := server.ListenAndServe(); err != nil && !apperrors.Is(err, http.ErrServerClosed) {
//			log.Printf("[LB] error for service=%s: %v", service, err)
//		}
//	}()
//}
//
///*
//StopServiceListener gracefully stops the load balancer
//for a given service.
//Call this when replicas reach 0 or service is deleted.
//*/
//func (lb *LoadBalancer) StopServiceListener(service string) {
//	lb.mu.Lock()
//	server, exists := lb.servers[service]
//	if !exists {
//		lb.mu.Unlock()
//		return
//	}
//	delete(lb.servers, service)
//	lb.mu.Unlock()
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	log.Printf("[LB] stopping service=%s", service)
//	_ = server.Shutdown(ctx)
//}
//
///*
//serviceHandler routes incoming requests to running pods
//using round-robin selection.
//*/
//func (lb *LoadBalancer) serviceHandler(service string) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		pods, err := lb.podRepo.ListPodsByService(context.Background(), service)
//		if err != nil {
//			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
//			return
//		}
//
//		runningPods := lb.filterRunningPods(pods)
//		if len(runningPods) == 0 {
//			http.Error(w, "no running pods", http.StatusServiceUnavailable)
//			return
//		}
//
//		target := lb.selectNextPod(runningPods)
//		log.Println("[LB] proxying to", target.String())
//
//		httputil.NewSingleHostReverseProxy(target).ServeHTTP(w, r)
//	}
//}
//
///*
//filterRunningPods returns only pods in Running state.
//*/
//func (lb *LoadBalancer) filterRunningPods(pods []entities.Pod) []entities.Pod {
//	out := make([]entities.Pod, 0)
//	for _, pod := range pods {
//		if pod.Status == entities.PodRunning {
//			out = append(out, pod)
//		}
//	}
//	return out
//}
//
///*
//selectNextPod selects a pod using round-robin strategy.
//*/
//func (lb *LoadBalancer) selectNextPod(pods []entities.Pod) *url.URL {
//	idx := int(atomic.AddUint64(&lb.counter, 1)) % len(pods)
//	return &url.URL{
//		Scheme: "http",
//		Host:   pods[idx].IP + ":80",
//	}
//}
