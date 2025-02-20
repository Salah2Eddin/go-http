package errors

import "fmt"

type ErrInvalidRequestLine struct{}

func (err ErrInvalidRequestLine) Error() string {
	return "Request is not an HTTP request"
}

type ErrInvalidHeader struct {
	Line string
}

func (err ErrInvalidHeader) Error() string {
	return fmt.Sprintf("%s is an invalid header", err.Line)
}
