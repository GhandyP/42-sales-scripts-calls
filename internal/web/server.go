package web

import (
	"net/http"
	"time"
)

const DefaultListenAddress = "127.0.0.1:8080"

// NewServer returns a loopback-only HTTP server with bounded request and
// connection timeouts. Callers provide the handler so server policy stays
// separate from route behavior.
func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              DefaultListenAddress,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
