package pkgerrors

type HTTPError interface {
	error
	HTTPStatusCode() int
}
