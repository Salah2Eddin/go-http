package pkgerrors

type ErrInvalidUri struct {
	Uri string
}

func (err *ErrInvalidUri) Error() string {
	return "Invalid URI: " + err.Uri
}
