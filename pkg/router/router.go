package router

import (
	"ducky/http/pkg/request"
	"ducky/http/pkg/response"
	"ducky/http/pkg/response/statuscodes"
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

func (router *Router) GetRoute(uri string) (*route, bool) {
	route, exists := router.routes[uri]
	return route, exists
}

func (router *Router) RouteRequest(request *request.Request) *response.Response {
	uri := request.Uri()
	route, exists := router.routes[uri]
	if !exists {
		return response.NewErrorResponse(statuscodes.Status404())
	}
	response := route.handle(request)
	return response
}
