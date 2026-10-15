package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServerBindsLoopbackAndSetsTimeouts(t *testing.T) {
	server := NewServer(http.NotFoundHandler())

	if got, want := server.Addr, "127.0.0.1:8080"; got != want {
		t.Fatalf("server address = %q, want loopback address %q", got, want)
	}
	if server.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("ReadHeaderTimeout = %s, want 5s", server.ReadHeaderTimeout)
	}
	if server.ReadTimeout != 30*time.Second {
		t.Errorf("ReadTimeout = %s, want 30s", server.ReadTimeout)
	}
	if server.WriteTimeout != 120*time.Second {
		t.Errorf("WriteTimeout = %s, want 120s", server.WriteTimeout)
	}
	if server.IdleTimeout != 60*time.Second {
		t.Errorf("IdleTimeout = %s, want 60s", server.IdleTimeout)
	}
	if server.MaxHeaderBytes != 1<<20 {
		t.Errorf("MaxHeaderBytes = %d, want %d", server.MaxHeaderBytes, 1<<20)
	}
}

func TestNewServerUsesProvidedHandler(t *testing.T) {
	server := NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("ready"))
	}))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	server.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusAccepted)
	}
	if got, want := recorder.Body.String(), "ready"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}
