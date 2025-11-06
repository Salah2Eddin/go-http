package pkgerrors

import "fmt"

type ErrInvalidRequestLine struct {
}

func (err ErrInvalidRequestLine) Error() string {
	return "Request is not an HTTP request"
}

func (err ErrInvalidRequestLine) HTTPStatusCode() int {
	return 400
}

type ErrInvalidHeader struct {
	Reason string
}

func (err ErrInvalidHeader) Error() string {
	return fmt.Sprintf("Invalid header line because %s", err.Reason)
}

func (err ErrInvalidHeader) HTTPStatusCode() int {
	return 400
}

type ErrInvalidContentLength struct {
	Length string
}

func (err ErrInvalidContentLength) Error() string {
	return fmt.Sprintf("%s is an invalid content length", err.Length)
}

func (err ErrInvalidContentLength) HTTPStatusCode() int {
	return 400
}

type ErrIncorrectContentLength struct {
	Cause error
}

func (err ErrIncorrectContentLength) Error() string {
	return fmt.Sprintf("Incorrect content length: %v", err.Cause)
}

func (err ErrIncorrectContentLength) HTTPStatusCode() int {
	return 400
}

type ErrExpectedEmptyBody struct{}

func (err ErrExpectedEmptyBody) Error() string {
	return fmt.Sprintf("Expected empty body")
}

func (err ErrExpectedEmptyBody) HTTPStatusCode() int {
	return 400
}

type ErrUnsupportedBodyTransferEncoding struct {
}

func (err ErrUnsupportedBodyTransferEncoding) Error() string {
	return fmt.Sprintf("Unsupported body transfer encoding")
}

func (err ErrUnsupportedBodyTransferEncoding) HTTPStatusCode() int {
	return 400
}
