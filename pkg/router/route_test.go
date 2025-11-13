package router

import (
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
)

func TestNewRoute(t *testing.T) {
	route := newRoute()
	if route == nil {
		t.Fatal("newRoute() returned nil")
	}
	if route.methodHandlers == nil {
		t.Errorf("newRoute() methodHandlers not initialized")
	}
	if len(route.methodHandlers) != 0 {
		t.Errorf("newRoute() methodHandlers not empty")
	}
}

func TestRouteAddHandler(t *testing.T) {
	route := newRoute()
	handler := func(req *request.Request) (*response.Response, error) {
		return nil, nil
	}
	route.AddHandler("GET", handler)
	if len(route.methodHandlers) != 1 {
		t.Errorf("AddHandler() did not add handler")
	}
}

func TestRouteGetHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		exists bool
	}{
		{"get handler", "GET", true},
		{"get missing handler", "POST", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := newRoute()
			handler := func(req *request.Request) (*response.Response, error) {
				return nil, nil
			}
			route.AddHandler("GET", handler)

			result := route.GetHandler(tc.method)
			if tc.exists && result == nil {
				t.Errorf("GetHandler(%q) returned nil, want handler", tc.method)
			}
			if !tc.exists && result != nil {
				t.Errorf("GetHandler(%q) returned handler, want nil", tc.method)
			}
		})
	}
}

func TestRouteAddHandlerMultiple(t *testing.T) {
	route := newRoute()
	handler1 := func(req *request.Request) (*response.Response, error) {
		return nil, nil
	}
	handler2 := func(req *request.Request) (*response.Response, error) {
		return nil, nil
	}

	route.AddHandler("GET", handler1)
	route.AddHandler("POST", handler2)

	if route.GetHandler("GET") == nil {
		t.Errorf("GetHandler(\"GET\") returned nil")
	}
	if route.GetHandler("POST") == nil {
		t.Errorf("GetHandler(\"POST\") returned nil")
	}
}
