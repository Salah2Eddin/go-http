package pkgerrors

type AppError struct {
	wrapped *error
}

// HTTPStatusCode returns the appropriate HTTP status code for each error
func (app AppError) HTTPStatusCode() int {
	switch (*app.wrapped).(type) {
	case ErrInvalidRequestLine, ErrInvalidHeader, ErrInvalidUri, ErrInvalidContentLength, ErrIncorrectContentLength:
		return 400 // Bad Request
	case ErrRouteNotFound:
		return 404 // Not Found
	default:
		return 500 // Internal Server Error
	}
}

func NewAppError(wrapped *error) AppError {
	return AppError{
		wrapped: wrapped,
	}
}

func (app AppError) Error() string {
	return (*app.wrapped).Error()
}
