package router

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/handlers"
	"ducky/http/pkg/request"
)

type Router struct {
	routes map[string]handlers.Handler
}

func (router *Router) AddRoute(uri string, handler handlers.Handler) {
	router.routes[uri] = handler
}

func (router *Router) Route(uri string, request request.Request) error {
	handler, exists := router.routes[uri]
	if !exists {
		return &errors.ErrInvalidRoute{Uri: uri}
	}
	err := handler.Handle(request)
	return err
}
