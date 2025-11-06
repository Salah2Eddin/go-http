package pkgerrors

import "fmt"

type ErrRouteExists struct {
	Route string
}

func (err ErrRouteExists) Error() string {
	return fmt.Sprintf("%s already exists", err.Route)
}

type ErrRouteNotFound struct {
	Route string
}

func (err ErrRouteNotFound) Error() string {
	return fmt.Sprintf("%s doesn't exist", err.Route)
}

func (err ErrRouteNotFound) HTTPStatusCode() int {
	return 404
}

type ErrMethodNotAllowed struct {
	Method string
}

func (err ErrMethodNotAllowed) Error() string {
	return fmt.Sprintf("Method %s is not allowed", err.Method)
}

func (err ErrMethodNotAllowed) HTTPStatusCode() int {
	return 405
}
