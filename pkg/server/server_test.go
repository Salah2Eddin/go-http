package server

import (
	"errors"
	"github.com/Salah2Eddin/go-http/pkg/readers"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/router"
	"github.com/Salah2Eddin/go-http/pkg/serializers"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer()

	if server == nil {
		t.Fatal("NewServer() returned nil, expected non-nil server instance")
	}

	if server.addr.Port != "8576" {
		t.Errorf("expected default port 8576, got %s", server.addr.Port)
	}

	if server.router == nil {
		t.Error("NewServer() returned server with nil router, expected non-nil router to be initialized")
	}

	if server.serializer == nil {
		t.Error("NewServer() returned server with nil serializer, expected non-nil serializer to be initialized")
	}

	if server.reader == nil {
		t.Error("NewServer() returned server with nil reader, expected non-nil reader to be initialized")
	}

	if server.errorResponder == nil {
		t.Error("NewServer() returned server with nil errorResponder, expected non-nil errorResponder to be initialized")
	}
}
func TestNewServerWithOptions(t *testing.T) {
	customPort := "9000"

	server := NewServer(WithPort(customPort))

	if server.addr.Port != customPort {
		t.Errorf("expected port %s, got %s", customPort, server.addr.Port)
	}
}

func TestNewServerMultipleOptions(t *testing.T) {
	customPort := "9000"
	customIP := "192.168.1.1"

	server := NewServer(WithPort(customPort), WithIP(customIP))

	if server.addr.Port != customPort {
		t.Errorf("expected port %s, got %s", customPort, server.addr.Port)
	}

	if server.addr.IP != customIP {
		t.Errorf("expected IP %s, got %s", customIP, server.addr.IP)
	}
}

func TestNewServerDefaultComponents(t *testing.T) {
	server := NewServer()

	if _, ok := server.router.(*router.Router); !ok {
		t.Errorf("expected default router type *router.Router, got %T", server.router)
	}

	if _, ok := server.serializer.(serializers.ResponseSerializer); !ok {
		t.Errorf("expected default serializer type *serializers.ResponseSerializer, got %T", server.serializer)
	}

	if _, ok := server.reader.(readers.RequestReader); !ok {
		t.Errorf("expected default reader type *readers.RequestReader, got %T", server.reader)
	}

	if _, ok := server.errorResponder.(JSONErrorResponder); !ok {
		t.Errorf("expected default errorResponder type JSONErrorResponder, got %T", server.errorResponder)
	}
}

func TestAddHandler(t *testing.T) {
	server := NewServer()

	handler := func(req *request.Request) (*response.Response, error) {
		return &response.Response{}, nil
	}

	err := server.AddHandler("/test", "GET", handler)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestAddHandlerMultiple(t *testing.T) {
	server := NewServer()

	handler1 := func(req *request.Request) (*response.Response, error) {
		return &response.Response{}, nil
	}

	handler2 := func(req *request.Request) (*response.Response, error) {
		return &response.Response{}, nil
	}

	err := server.AddHandler("/users", "GET", handler1)
	if err != nil {
		t.Errorf("expected no error for first handler, got %v", err)
	}

	err = server.AddHandler("/posts", "POST", handler2)
	if err != nil {
		t.Errorf("expected no error for second handler, got %v", err)
	}
}

func TestCloseListener(t *testing.T) {
	listener := &mockListener{}

	closeListener(listener)

	if !listener.closed {
		t.Error("listener.closed is false, expected listener to be closed after closeListener()")
	}
}

func TestCloseListenerWithError(t *testing.T) {
	listener := &mockListener{
		closeError: errors.New("close error"),
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on close error")
		}
	}()

	closeListener(listener)
}

// TODO: test server start and accept loop
