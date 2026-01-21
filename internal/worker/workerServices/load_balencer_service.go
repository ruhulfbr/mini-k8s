package workerServices

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type LoadBalancer struct {
	mu       sync.RWMutex
	clusters map[int64]*clusterState
}

type clusterState struct {
	clusterId int64
	port      int
	server    *http.Server

	counter uint64
	pods    map[int64]*entities.Pod // podId → pod
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		clusters: make(map[int64]*clusterState),
	}
}

func (lb *LoadBalancer) AddPod(clusterId int64, port int, pod *entities.Pod) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	state, exists := lb.clusters[clusterId]
	if !exists {
		state = &clusterState{
			clusterId: clusterId,
			port:      port,
			pods:      make(map[int64]*entities.Pod),
		}
		lb.clusters[clusterId] = state
		lb.startClusterLocked(state)
	}

	state.pods[pod.Id] = pod

	log.Printf(
		"[LB] pod added cluster=%d pod=%d (%s)",
		clusterId, pod.Id, pod.IpAddress,
	)
}

func (lb *LoadBalancer) RemovePod(clusterId int64, podId int64) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	state, exists := lb.clusters[clusterId]
	if !exists {
		return
	}

	delete(state.pods, podId)

	log.Printf(
		"[LB] pod removed cluster=%d pod=%d",
		clusterId, podId,
	)

	if len(state.pods) == 0 {
		lb.stopClusterLocked(state)
		delete(lb.clusters, clusterId)
	}
}

func (lb *LoadBalancer) startClusterLocked(state *clusterState) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", state.port),
		Handler: lb.clusterHandler(state.clusterId),
	}

	state.server = server

	go func() {
		log.Printf(
			"[LB] started cluster=%d port=%d",
			state.clusterId, state.port,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Printf(
				"[LB] error cluster=%d: %v",
				state.clusterId, err,
			)
		}
	}()
}

func (lb *LoadBalancer) stopClusterLocked(state *clusterState) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("[LB] stopping cluster=%d", state.clusterId)
	_ = state.server.Shutdown(ctx)
}

func (lb *LoadBalancer) clusterHandler(clusterId int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lb.mu.RLock()
		state, exists := lb.clusters[clusterId]
		if !exists || len(state.pods) == 0 {
			lb.mu.RUnlock()
			http.Error(w, "cluster unavailable", http.StatusServiceUnavailable)
			return
		}

		pods := make([]*entities.Pod, 0, len(state.pods))
		for _, pod := range state.pods {
			pods = append(pods, pod)
		}
		lb.mu.RUnlock()

		idx := int(atomic.AddUint64(&state.counter, 1)) % len(pods)
		target := pods[idx]

		targetURL := &url.URL{
			Scheme: "http",
			Host:   target.IpAddress + ":80",
		}

		log.Printf(
			"[LB] cluster=%d -> pod=%d (%s)",
			clusterId, target.Id, target.IpAddress,
		)

		httputil.NewSingleHostReverseProxy(targetURL).
			ServeHTTP(w, r)
	}
}
