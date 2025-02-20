package request

type RequestLine struct {
	Method  string
	Uri     string
	Version string
}

func NewRequestLine(method string, uri string, version string) *RequestLine {
	return &RequestLine{
		Method:  method,
		Uri:     uri,
		Version: version,
	}
}
