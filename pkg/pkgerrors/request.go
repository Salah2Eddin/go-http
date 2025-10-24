package pkgerrors

import "fmt"

type ErrInvalidRequestLine struct {
}

func (err ErrInvalidRequestLine) Error() string {
	return "Request is not an HTTP request"
}

type ErrInvalidHeader struct {
	Reason string
}

func (err ErrInvalidHeader) Error() string {
	return fmt.Sprintf("Invalid header line because %s", err.Reason)
}

type ErrInvalidContentLength struct {
	Length string
}

func (err ErrInvalidContentLength) Error() string {
	return fmt.Sprintf("%s is an invalid content length", err.Length)
}

type ErrIncorrectContentLength struct {
	Cause error
}

func (err ErrIncorrectContentLength) Error() string {
	return fmt.Sprintf("Incorrect content length: %v", err.Cause)
}

type ErrExpectedEmptyBody struct{}

func (e ErrExpectedEmptyBody) Error() string {
	return fmt.Sprintf("Expected empty body")
}

type ErrUnsupportedBodyTransferEncoding struct {
}

func (e ErrUnsupportedBodyTransferEncoding) Error() string {
	return fmt.Sprintf("Unsupported body transfer encoding")
}
