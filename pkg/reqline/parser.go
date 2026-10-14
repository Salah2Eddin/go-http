package reqline

import (
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/uri"
	"github.com/Salah2Eddin/go-http/pkg/util/charutil"
	"strings"
)

func validRequestLine(parts []string) bool {
	if len(parts) != 3 {
		return false
	}
	httpVer := parts[2]
	return strings.HasPrefix(httpVer, "HTTP/")
}

func ParseRequestLine(requestLineBytes []byte) (*RequestLine, error) {

	// Request line must contain bytes in the ASCII range only (RFC9112 2.2)
	if !charutil.ValidateAsciiEncoding(requestLineBytes) {
		return nil, pkgerrors.ErrInvalidRequestLine{}
	}

	requestLine := string(requestLineBytes)
	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Fields(requestLine)

	if !validRequestLine(parts) {
		return nil, pkgerrors.ErrInvalidRequestLine{}
	}

	method := parts[0]
	uriString := parts[1]
	httpVer := parts[2]

	if !uri.ValidateURI(uriString) {
		return nil, pkgerrors.ErrInvalidUri{Uri: uriString}
	}

	line := NewRequestLine(
		method,                // method
		uri.NewUri(uriString), // path
		httpVer,               // http version
	)
	return line, nil
}
