package router

import (
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/response/statuscodes"
)

type Handler func(request request.Request) response.Response

type Route struct {
	methodHandlers map[string]Handler
}

func newRoute() Route {
	route := Route{}
	route.methodHandlers = make(map[string]Handler)
	return route
}

func (route *Route) AddHandler(method string, handler Handler) {
	route.methodHandlers[method] = handler
}

func (route *Route) handle(request request.Request) response.Response {
	method := request.Method()
	handler, exists := route.methodHandlers[method]
	if !exists {
		// should be Error 405 but it wasn't introduced in HTTP/1.0
		// so i settled with error 400 instead
		return response.NewEmptyResponse(statuscodes.Status400())
	}

	resp := handler(request)
	return resp
}
