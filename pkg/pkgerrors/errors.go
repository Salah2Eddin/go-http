package pkgerrors

type AppError struct {
	wrapped error
}

func NewAppError(wrapped error) *AppError {
	return &AppError{
		wrapped: wrapped,
	}
}

// HTTPStatusCode returns the appropriate HTTP status code for each error
func (app AppError) HTTPStatusCode() string {
	switch app.wrapped.(type) {
	case ErrInvalidRequestLine, ErrInvalidHeader, ErrInvalidUri, ErrInvalidContentLength, ErrIncorrectContentLength, ErrExpectedEmptyBody, ErrUnsupportedBodyTransferEncoding:
		return "400" // Bad Request
	case ErrMethodNotAllowed:
		return "405" // Method Not Allowed
	case ErrRouteNotFound:
		return "404" // Not Found
	default:
		return "500" // Internal Server Error
	}
}

func (app AppError) Error() string {
	return app.wrapped.Error()
}
