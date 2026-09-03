package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Route struct {
	Prefix string
	Target string
}

type Gateway struct {
	routes  []Route
	proxies map[string]*httputil.ReverseProxy
	health  map[string]bool
	mu      sync.RWMutex
	sem     chan struct{}
}

func NewGateway(routes []Route, maxConcurrent int) *Gateway {
	proxies := make(map[string]*httputil.ReverseProxy)
	health := make(map[string]bool)

	for _, r := range routes {
		target, err := url.Parse(r.Target)
		if err != nil {
			log.Fatalf("invalid target url %s: %v", r.Target, err)
		}
		proxies[r.Prefix] = httputil.NewSingleHostReverseProxy(target)
		health[r.Target] = false
	}

	return &Gateway{
		routes:  routes,
		proxies: proxies,
		health:  health,
		sem:     make(chan struct{}, maxConcurrent),
	}
}

func (gw *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range gw.routes {
		if strings.HasPrefix(r.URL.Path, route.Prefix) {
			gw.sem <- struct{}{}
			defer func() { <-gw.sem }()

			log.Printf("%s %s -> %s", r.Method, r.URL.Path, route.Target)
			gw.proxies[route.Prefix].ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (gw *Gateway) startHeartbeatMonitor(interval time.Duration) {
	go func() {
		for {
			gw.checkHealth()
			time.Sleep(interval)
		}
	}()
}

func (gw *Gateway) checkHealth() {
	for _, route := range gw.routes {
		resp, err := http.Get(route.Target + "/health")
		healthy := err == nil && resp.StatusCode == http.StatusOK
		if resp != nil {
			resp.Body.Close()
		}

		gw.mu.Lock()
		gw.health[route.Target] = healthy
		gw.mu.Unlock()
	}
}

func (gw *Gateway) statusHandler(w http.ResponseWriter, r *http.Request) {
	gw.mu.RLock()
	defer gw.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gw.health)
}
