package pkgerrors

import "errors"

type AppError struct {
	wrapped error
	code    int
}

func NewAppError(wrapped error) *AppError {
	return &AppError{
		wrapped: wrapped,
	}
}

func (err AppError) HTTPStatusCode() int {
	var httpErr HTTPError
	if errors.As(err.wrapped, &httpErr) {
		return httpErr.HTTPStatusCode()
	}
	return 500
}

func (err AppError) Error() string {
	return err.wrapped.Error()
}
