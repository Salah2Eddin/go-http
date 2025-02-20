package router

import (
	"ducky/http/pkg/request"
	"ducky/http/pkg/response"
	statuscodes "ducky/http/pkg/response/status_codes"
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

func (router *Router) Route(request *request.Request) *response.Response {
	uri := request.Line.Uri
	route, exists := router.routes[uri]
	if !exists {
		return response.NewErrorResponse(statuscodes.Status404())
	}
	response := route.handle(request)
	return response
}
