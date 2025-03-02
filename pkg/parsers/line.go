package parsers

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
	"ducky/http/pkg/uri"
	"ducky/http/pkg/util"
	"strings"
)

func validRequestLine(parts []string) bool {
	if len(parts) != 3 {
		return false
	}

	http_ver := parts[2]

	return strings.HasPrefix(http_ver, "HTTP")
}

func validateAsciiEncoding(bytes *[]byte) bool {
	for _, v := range *bytes {
		if !util.IsUSASCII(v) {
			return false
		}
	}
	return true
}

func parseRequestLine(request_line_bytes *[]byte) (*request.RequestLine, error) {

	// Request line must contain bytes in the USASCII range only (RFC9112 2.2)
	if !validateAsciiEncoding(request_line_bytes) {
		return &request.RequestLine{}, errors.ErrInvalidRequestLine{}
	}

	request_line := string(*request_line_bytes)
	request_line = strings.TrimSpace(request_line)
	parts := strings.Fields(request_line)

	if !validRequestLine(parts) {
		return &request.RequestLine{}, errors.ErrInvalidRequestLine{}
	}

	uri := uri.NewUri(parts[1])

	return request.NewRequestLine(
		parts[0], // method
		uri,
		parts[2], // http version
	), nil
}
