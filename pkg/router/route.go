package router

import (
	"ducky/http/pkg/handlers"
	"ducky/http/pkg/request"
)

type Route struct {
	method_handlers map[string]handlers.Handler
}

func (route *Route) AddHandler(method string, handler handlers.Handler) {
	route.method_handlers[method] = handler
}

func (route *Route) Handle(request *request.Request) error {
	method := request.Line.Method
	handler, exists := route.method_handlers[method]
	if !exists {
		return &errors.ErrMethodNotAllowed{
			Method: request.Line.Method,
			Uri:    request.Line.Uri,
		}
	}

	err := handler.Handle(request)
	return err
}
