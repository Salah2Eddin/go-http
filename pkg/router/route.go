package router

import (
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
)

type Handler func(request *request.Request) (*response.Response, pkgerrors.HTTPError)

type Route struct {
	methodHandlers map[string]Handler
}

func newRoute() *Route {
	route := Route{}
	route.methodHandlers = make(map[string]Handler)
	return &route
}

func (route *Route) AddHandler(method string, handler Handler) {
	route.methodHandlers[method] = handler
}

func (route *Route) GetHandler(method string) Handler {
	return route.methodHandlers[method]
}
