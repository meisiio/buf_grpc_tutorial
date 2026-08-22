package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	identityURL, _ := url.Parse("http://localhost:8081")
	identityProxy := httputil.NewSingleHostReverseProxy(identityURL)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/identity.") {
			log.Printf("[GATEWAY] Forwarding request to Identity Service: %s", r.URL.Path)
			identityProxy.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Microservice not found", http.StatusNotFound)
	})

	log.Println("API Gateway listening on public port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Gateway failed: %v", err)
	}
}
