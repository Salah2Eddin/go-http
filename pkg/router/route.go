package router

import (
	"ducky/http/pkg/request"
	"ducky/http/pkg/response"
	"ducky/http/pkg/response/statuscodes"
)

type Handler func(request *request.Request) *response.Response

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

func (route *route) handle(request *request.Request) *response.Response {
	method := request.Line.Method
	handler, exists := route.method_handlers[method]
	if !exists {
		// should be Error 405 but it wasn't introduced in HTTP/1.0
		// so i settled with error 400 instead
		return response.NewErrorResponse(statuscodes.Status400())
	}

	response := handler(request)
	return response
}
