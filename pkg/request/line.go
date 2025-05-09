package request

import "github.com/Salah2Eddin/go-http/pkg/uri"

type RequestLine struct {
	Method  string
	Uri     *uri.Uri
	Version string
}

func NewRequestLine(method string, uri *uri.Uri, version string) *RequestLine {
	return &RequestLine{
		Method:  method,
		Uri:     uri,
		Version: version,
	}
}
