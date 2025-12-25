package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type LoadBalancer struct {
	ds       *datastore.Datastore
	counter  uint64
	interval time.Duration
}

func NewLoadBalancer(ds *datastore.Datastore) *LoadBalancer {
	return &LoadBalancer{ds: ds, interval: 5 * time.Second}
}

func (lb *LoadBalancer) Serve() {

	nodes, err := lb.ds.ListPodsGroupedByService(context.Background())

	if err != nil || len(nodes) == 0 {
		log.Println("No Nodes available to run")
		return
	}

	fmt.Println("Nodes:", nodes)

	for _, pod := range nodes {
		go lb.RunForNode(pod.Service, pod.Port)
	}
}

func (lb *LoadBalancer) RunForNode(service string, port int) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pods, err := lb.ds.ListPodsByService(context.Background(), service)
		if err != nil || len(pods) == 0 {
			http.Error(w, "no pods available", http.StatusServiceUnavailable)
			return
		}

		var runningPods []entities.Pod
		for _, pod := range pods {
			if pod.Status == "Running" {
				runningPods = append(runningPods, pod)
			}
		}

		if len(runningPods) == 0 {
			http.Error(w, "no running pods", http.StatusServiceUnavailable)
			return
		}

		idx := int(atomic.AddUint64(&lb.counter, 1)) % len(runningPods)
		target := &url.URL{
			Scheme: "http",
			Host:   runningPods[idx].IP + ":80",
		}

		fmt.Println("Load balancer target to : ", target.String())

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
	}

	log.Printf("Load balancer listening on port %d", port)
	log.Fatal(server.ListenAndServe())
}

//func (lb *LoadBalancer) runForNodes(pod entities.Pod) {
//	handler := func(w http.ResponseWriter, r *http.Request) {
//		pods, err := lb.ds.ListPodsByService(context.Background(), pod.Service)
//		if err != nil || len(pods) == 0 {
//			http.Error(w, "no pods available", http.StatusServiceUnavailable)
//			return
//		}
//
//		var runningPods []entities.Pod
//		for _, pod := range pods {
//			if pod.Status == "Running" {
//				runningPods = append(runningPods, pod)
//			}
//		}
//
//		if len(runningPods) == 0 {
//			http.Error(w, "no running pods", http.StatusServiceUnavailable)
//			return
//		}
//
//		idx := int(atomic.AddUint64(&lb.counter, 1)) % len(runningPods)
//
//		target := &url.URL{
//			Scheme: "http",
//			Host:   runningPods[idx].IP + ":80",
//		}
//
//		fmt.Println("Load balancer target to : ", target.String())
//
//		proxy := httputil.NewSingleHostReverseProxy(target)
//		proxy.ServeHTTP(w, r)
//	}
//
//	server := &http.Server{
//		Addr:    fmt.Sprintf(":%d", pod.Port),
//		Handler: http.HandlerFunc(handler),
//	}
//
//	log.Printf("Load balancer listening on port %d", pod.Port)
//	log.Fatal(server.ListenAndServe())
//}
