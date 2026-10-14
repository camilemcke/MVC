
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "sync"
    "time"
)

type ConsulService struct {
    Service struct {
        Address string `json:"Address"`
        Port    int    `json:"Port"`
    } `json:"Service"`
}

type LoadBalancer struct {
    current int
    mu      sync.Mutex
    client  *http.Client
}

func (lb *LoadBalancer) getHealthyServers() ([]string, error) {
    consulAddress := os.Getenv("CONSUL_HTTP_ADDR")
    if consulAddress == "" {
        consulAddress = "http://consul:8500"
    }

    endpoint := consulAddress + "/v1/health/service/products?passing=true"

    response, err := lb.client.Get(endpoint)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Consul returned status %d", response.StatusCode)
    }

    var services []ConsulService

    if err := json.NewDecoder(response.Body).Decode(&services); err != nil {
        return nil, err
    }

    servers := make([]string, 0, len(services))

    for _, service := range services {
        if service.Service.Address == "" || service.Service.Port == 0 {
            continue
        }

        address := fmt.Sprintf(
            "http://%s:%d",
            service.Service.Address,
            service.Service.Port,
        )

        servers = append(servers, address)
    }

    return servers, nil
}

func (lb *LoadBalancer) nextServer(servers []string) string {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    server := servers[lb.current%len(servers)]
    lb.current++

    return server
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/health" {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("load balancer alive"))
        return
    }

    servers, err := lb.getHealthyServers()
    if err != nil {
        log.Println("Error querying Consul:", err)
        http.Error(w, "service discovery unavailable", http.StatusServiceUnavailable)
        return
    }

    if len(servers) == 0 {
        http.Error(w, "no healthy product servers available", http.StatusServiceUnavailable)
        return
    }

    targetURL := lb.nextServer(servers)

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
        client: &http.Client{
            Timeout: 3 * time.Second,
        },
    }

    log.Println("load balancer listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", lb))
}
