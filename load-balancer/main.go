package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	servers []string
	current int
	mu      sync.Mutex
}

func (lb *LoadBalancer) nextServer() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.current]
	lb.current = (lb.current + 1) % len(lb.servers)

	return server
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("load balancer alive"))
		return
	}
	targetURL := lb.nextServer()

	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "invalid target", http.StatusInternalServerError)
		return
	}

	log.Printf("%s %s -> %s", r.Method, r.URL.Path, targetURL)

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

func main() {
	lb := &LoadBalancer{
		servers: []string{
			"http://backend-1:8080",
			"http://backend-2:8080",
		},
	}

	log.Println("load balancer listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", lb))
}