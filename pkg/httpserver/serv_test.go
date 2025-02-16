package httpserver

import (
	"net"
	"net/http"
	"testing"
	"time"
)

func TestPort(t *testing.T) {
	s := &Server{
		server: &http.Server{},
	}
	Port("8081")(s)

	expectedAddr := net.JoinHostPort("", "8081")
	if s.server.Addr != expectedAddr {
		t.Errorf("expected server address to be %s, got %s", expectedAddr, s.server.Addr)
	}
}

func TestReadTimeout(t *testing.T) {
	s := &Server{
		server: &http.Server{},
	}
	timeout := 10 * time.Second
	ReadTimeout(timeout)(s)

	if s.server.ReadTimeout != timeout {
		t.Errorf("expected read timeout to be %v, got %v", timeout, s.server.ReadTimeout)
	}
}

func TestWriteTimeout(t *testing.T) {
	s := &Server{
		server: &http.Server{},
	}
	timeout := 10 * time.Second
	WriteTimeout(timeout)(s)

	if s.server.WriteTimeout != timeout {
		t.Errorf("expected write timeout to be %v, got %v", timeout, s.server.WriteTimeout)
	}
}

func TestShutdownTimeout(t *testing.T) {
	s := &Server{
		server: &http.Server{},
	}
	timeout := 10 * time.Second
	ShutdownTimeout(timeout)(s)

	if s.shutdownTimeout != timeout {
		t.Errorf("expected shutdown timeout to be %v, got %v", timeout, s.shutdownTimeout)
	}
}
