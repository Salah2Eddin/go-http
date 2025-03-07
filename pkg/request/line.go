package request

import "ducky/http/pkg/uri"

type Line struct {
	Method  string
	Uri     uri.Uri
	Version string
}

func NewRequestLine(method string, uri uri.Uri, version string) Line {
	return Line{
		Method:  method,
		Uri:     uri,
		Version: version,
	}
}
