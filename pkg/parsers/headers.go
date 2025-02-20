package parsers

import (
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
	"strings"
)

// func validateHeaderLine(parts []string) bool {
// 	return len(parts) == 2
// }

func parseHeaderLine(line string) (string, string, error) {
	key, value, found := strings.Cut(line, ":")

	if !found {
		return "", "", errors.ErrInvalidHeader{Line: line}
	}

	// lower-case to make headers case-insensitive as specified in RFC1945 (HTTP 1.0 Specifications)
	key = strings.ToLower(key)

	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	return key, value, nil
}

func ParseRequestHeaders(lines []string) (*request.RequestHeaders, error) {
	headers := request.NewRequestHeaders()
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		key, value, err := parseHeaderLine(line)
		if err != nil {
			return &request.RequestHeaders{}, err
		}

		headers.Set(key, value)
	}

	return headers, nil
}
