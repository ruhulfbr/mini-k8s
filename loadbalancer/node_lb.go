package loadbalancer

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/ruhulfbr/mini-k8s/orchestrator"
)

type Manager struct {
	ds      orchestrator.Datastore
	mu      sync.Mutex
	active  map[string]bool // service -> running?
	counter uint64
}

func NewManager(ds orchestrator.Datastore) *Manager {
	return &Manager{
		ds:     ds,
		active: make(map[string]bool),
	}
}

func (m *Manager) StartNode(service string, port string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		pods, err := m.ds.ListRunningPods(ctx, service)
		if err != nil || len(pods) == 0 {
			http.Error(w, "No pods available", http.StatusServiceUnavailable)
			return
		}

		// round-robin
		idx := atomic.AddUint64(&m.counter, 1)
		pod := pods[int(idx)%len(pods)]

		target := &url.URL{
			Scheme: "http",
			Host:   pod.IP + ":" + string(rune(pod.Port)),
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	})

	log.Printf("[NodeLB:%s] listening on :%s\n", service, port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
