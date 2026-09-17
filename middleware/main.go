package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	routes := []Route{
		{Prefix: "/user", Target: "http://backend-1:8080"},
		{Prefix: "/product", Target: "http://load-balancer:8080"},
		{Prefix: "/inventory", Target: "http://inventory:8080"},
		{Prefix: "/analytics", Target: "http://analytics:8080"},
	}

	gw := NewGateway(routes, 10)
	gw.startHeartbeatMonitor(5 * time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/heartbeat", heartbeatHandler)
	mux.HandleFunc("/status", gw.statusHandler)
	mux.Handle("/", gw)

	log.Println("gateway listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", cors(mux)))
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("gateway alive"))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
