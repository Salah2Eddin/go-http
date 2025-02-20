package router

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
)

type Router struct {
	routes map[string]*route
}

func NewRouter() *Router {
	router := Router{}
	router.routes = make(map[string]*route)
	return &router
}

func (router *Router) AddRoute(uri string, route *route) {
	router.routes[uri] = route
}

func (router *Router) NewRoute(uri string) *route {
	route := newRoute()
	router.routes[uri] = route
	return route
}

func (router *Router) Route(request *request.Request) error {
	uri := request.Line.Uri
	route, exists := router.routes[uri]
	if !exists {
		return &errors.ErrInvalidRoute{Uri: uri}
	}
	err := route.handle(request)
	return err
}
