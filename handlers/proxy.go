package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// CreateProxyHandler creates a new reverse proxy handler.
func CreateProxyHandler(serviceURL string) http.HandlerFunc {
	target, err := url.Parse("http://" + serviceURL)
	if err != nil {
		log.Printf("Error parsing service URL %s: %v", serviceURL, err)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = target.Host
		proxy.ServeHTTP(w, r)
	}
}
