package router

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
)

type Handler func(request *request.Request) error

type route struct {
	method_handlers map[string]Handler
}

func newRoute() *route {
	route := route{}
	route.method_handlers = make(map[string]Handler)
	return &route
}

func (route *route) AddHandler(method string, handler Handler) {
	route.method_handlers[method] = handler
}

func (route *route) handle(request *request.Request) error {
	method := request.Line.Method
	handler, exists := route.method_handlers[method]
	if !exists {
		return &errors.ErrMethodNotAllowed{
			Method: request.Line.Method,
			Uri:    request.Line.Uri,
		}
	}

	err := handler(request)
	return err
}
