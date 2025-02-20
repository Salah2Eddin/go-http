package parsers

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
	"strings"
)

func validRequestLine(parts []string) bool {
	if len(parts) != 3 {
		return false
	}

	http_ver := parts[2]

	return strings.HasPrefix(http_ver, "HTTP")
}

func requestLineFromParts(parts []string) *request.RequestLine {
	request_line := request.NewRequestLine(
		parts[0],
		parts[1],
		parts[2],
	)
	return request_line
}

func ParseRequestLine(request_line string) (*request.RequestLine, error) {
	request_line = strings.TrimSpace(request_line)
	parts := strings.Fields(request_line)

	if !validRequestLine(parts) {
		return &request.RequestLine{}, errors.ErrInvalidRequestLine{}
	}

	return requestLineFromParts(parts), nil
}
