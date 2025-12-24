package orchestrator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type ServiceConfig struct {
	ServiceName string `json:"service_name"`
	Image       string `json:"image"`
	Replicas    int    `json:"replicas"`
	Port        int    `json:"port"`
}

type ClusterConfig struct {
	Services []ServiceConfig `json:"services"`
	mu       sync.RWMutex
}

func LoadConfig(path string) *ClusterConfig {
	cfg := &ClusterConfig{}
	go func() {
		for {
			data, err := ioutil.ReadFile(path)

			if err != nil {
				log.Printf("failed to read cluster.json: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			var tmp ClusterConfig
			if err := json.Unmarshal(data, &tmp); err != nil {
				log.Printf("failed to unmarshal cluster.json: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			cfg.mu.Lock()
			cfg.Services = tmp.Services
			cfg.mu.Unlock()
			time.Sleep(5 * time.Second)
		}
	}()

	return cfg
}

func (c *ClusterConfig) GetServices() []ServiceConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Services
}
