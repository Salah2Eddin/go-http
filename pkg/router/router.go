package router

import (
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/uri"
)

type Router struct {
	routes map[int]*Route
	tree   RoutesTree
}

func NewRouter() Router {
	router := Router{
		routes: make(map[int]*Route),
		tree:   NewRoutesTree(),
	}
	return router
}

func (router *Router) newRoute(uri *uri.Uri) (*Route, error) {
	id, err := router.tree.addRoute(uri)
	if err != nil {
		return nil, err
	}
	router.routes[id] = newRoute()
	return router.routes[id], nil
}

func (router *Router) getOrCreateRoute(uri *uri.Uri) (*Route, error) {
	route, err := router.getRoute(uri, false)
	if err != nil {
		route, err = router.newRoute(uri)
		if err != nil {
			return nil, err
		}
	}
	return route, nil
}

func (router *Router) getRoute(uri *uri.Uri, allowWildcardInURI bool) (*Route, error) {
	id, err := router.tree.find(uri, allowWildcardInURI)
	if err != nil {
		return nil, err
	}
	// at this point, a route with id is guaranteed to exist
	return router.routes[id], nil
}

// AddHandler Registers a new handler for the given URI and HTTP method.
// If the route corresponding to the URI does not exist, a new route is created.
func (router *Router) AddHandler(uri *uri.Uri, method string, handler Handler) error {
	route, err := router.getOrCreateRoute(uri)
	if err != nil {
		return err
	}
	route.AddHandler(method, handler)
	return nil
}

// GetRequestHandler takes an uri and a method and returns the handler associated with them
func (router *Router) GetRequestHandler(uri *uri.Uri, method string) (Handler, error) {
	route, err := router.getRoute(uri, true)
	if err != nil {
		return nil, err
	}
	handler := route.GetHandler(method)
	if handler == nil {
		return nil, pkgerrors.ErrMethodNotAllowed{Method: method, Uri: uri.String()}
	}
	return handler, nil
}
