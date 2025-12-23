package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ruhulfbr/mini-k8s/orchestrator"
)

type LoadBalancer struct {
	ds       *orchestrator.Datastore
	counter  uint64
	interval time.Duration
}

func NewLoadBalancer(ds *orchestrator.Datastore) *LoadBalancer {
	return &LoadBalancer{ds: ds, interval: 5 * time.Second}
}

func (lb *LoadBalancer) Serve(port string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pods, err := lb.ds.ListPodsByService(context.Background(), "")
		if err != nil || len(pods) == 0 {
			http.Error(w, "no pods available", http.StatusServiceUnavailable)
			return
		}

		var runningPods []orchestrator.Pod
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
			Host:   runningPods[idx].IP + ":" + strconv.Itoa(runningPods[idx].Port),
		}

		fmt.Println("Load balancer target to : ", target.String())

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.HandlerFunc(handler),
	}
	log.Printf("Load balancer listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}
