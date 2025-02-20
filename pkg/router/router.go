package router

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
)

type Router struct {
	routes map[string]Route
}

func (router *Router) AddRoute(uri string, route Route) {
	router.routes[uri] = route
}

func (router *Router) Route(request *request.Request) error {
	uri := request.Line.Uri
	route, exists := router.routes[uri]
	if !exists {
		return &errors.ErrInvalidRoute{Uri: uri}
	}
	err := route.Handle(request)
	return err
}
