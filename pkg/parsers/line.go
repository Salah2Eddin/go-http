package parsers

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
	"ducky/http/pkg/uri"
	"strings"
)

func validRequestLine(parts []string) bool {
	if len(parts) != 3 {
		return false
	}

	http_ver := parts[2]

	return strings.HasPrefix(http_ver, "HTTP")
}

func ParseRequestLine(request_line string) (*request.RequestLine, error) {
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
