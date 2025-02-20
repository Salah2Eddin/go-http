package errors

import "fmt"

type ErrInvalidRoute struct {
	Uri string
}

func (err ErrInvalidRoute) Error() string {
	return fmt.Sprintf("%s doesn't have a handler", err.Uri)
}

type ErrMethodNotAllowed struct {
	Method string
	Uri    string
}

func (err ErrMethodNotAllowed) Error() string {
	return fmt.Sprintf("%s %s is not implemented", err.Method, err.Uri)
}
