package router

import (
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/response/statuscodes"
	"github.com/Salah2Eddin/go-http/pkg/uri"
)

type Router struct {
	routes map[int]Route
	tree   RoutesTree
}

func NewRouter() Router {
	router := Router{
		routes: make(map[int]Route),
		tree:   NewRoutesTree(),
	}
	return router
}

func (router *Router) NewRoute(uri uri.Uri) (Route, error) {
	route := newRoute()
	id, err := router.tree.addRoute(uri)
	if err != nil {
		return Route{}, err
	}
	router.routes[id] = route
	return route, nil
}

func (router *Router) GetRoute(uri uri.Uri, allowWildcard bool) (Route, error) {
	id, err := router.tree.find(uri, allowWildcard)
	if err != nil {
		return Route{}, err
	}
	// at this point, a route with id is guaranteed to exist
	return router.routes[id], nil
}

func (router *Router) RouteRequest(request request.Request) response.Response {
	route, err := router.GetRoute(request.Uri(), true)
	if err != nil {
		return response.NewEmptyResponse(statuscodes.Status404())
	}
	resp := route.handle(request)
	return resp
}
