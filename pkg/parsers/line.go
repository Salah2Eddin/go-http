package parsers

import (
	"ducky/http/pkg/pkgerrors"
	"ducky/http/pkg/request"
	"ducky/http/pkg/uri"
	"ducky/http/pkg/util/charutil"
	"strings"
)

func validRequestLine(parts []string) bool {
	if len(parts) != 3 {
		return false
	}

	httpVer := parts[2]

	return strings.HasPrefix(httpVer, "HTTP")
}

func validateAsciiEncoding(bytes []byte) bool {
	for _, v := range bytes {
		if !charutil.IsASCII(v) {
			return false
		}
	}
	return true
}

func parseRequestLine(requestLineBytes []byte) (request.Line, error) {

	// Request line must contain bytes in the ASCII range only (RFC9112 2.2)
	if !validateAsciiEncoding(requestLineBytes) {
		return request.Line{}, pkgerrors.ErrInvalidRequestLine{}
	}

	requestLine := string(requestLineBytes)
	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Fields(requestLine)

	if !validRequestLine(parts) {
		return request.Line{}, pkgerrors.ErrInvalidRequestLine{}
	}

	method := parts[0]
	uriString := parts[1]
	httpVer := parts[2]

	if !uri.ValidateURI(uriString) {
		return request.Line{}, &pkgerrors.ErrInvalidUri{Uri: uriString}
	}

	uriObj := uri.NewUri(uriString)
	return request.NewRequestLine(
		method, // method
		uriObj,
		httpVer, // http version
	), nil
}
