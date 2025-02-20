package errors

import "fmt"

type ErrInvalidRoute struct {
	Uri string
}

func (err ErrInvalidRoute) Error() string {
	return fmt.Sprintf("%s doesn't have a handler", err.Uri)
}
