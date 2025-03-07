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
