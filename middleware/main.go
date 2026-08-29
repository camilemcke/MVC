package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	backend, err := url.Parse("http://backend:8080")
	if err != nil {
		log.Fatal(err)
	}

	routes := map[string]*httputil.ReverseProxy{
		"/user": httputil.NewSingleHostReverseProxy(backend),
	}

	http.HandleFunc("/", cors(func(w http.ResponseWriter, r *http.Request) {
		for prefix, worker := range routes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				log.Printf("%s %s -> backend", r.Method, r.URL.Path)
				worker.ServeHTTP(w, r)
				return
			}
		}
		log.Printf("%s %s -> sin worker", r.Method, r.URL.Path)
		http.NotFound(w, r)
	}))

	log.Println("Middleware listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r)
	}
}
