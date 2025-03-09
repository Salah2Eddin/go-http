package pkgerrors

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

type ErrInvalidContentLength struct {
	Length string
}

func (err ErrInvalidContentLength) Error() string {
	return fmt.Sprintf("Invalid content length: %s", err.Length)
}

type ErrIncorrectContentLength struct {
	Length int
}

func (err ErrIncorrectContentLength) Error() string {
	return fmt.Sprintf("Request body's length is not: %d", err.Length)
}
