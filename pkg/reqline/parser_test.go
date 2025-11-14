package reqline

import (
	"errors"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"testing"
)

func TestValidRequestLine(t *testing.T) {
	testCases := []struct {
		name     string
		parts    []string
		expected bool
	}{
		{
			name:     "valid three parts with HTTP prefix",
			parts:    []string{"GET", "/path", "HTTP/1.1"},
			expected: true,
		},
		{
			name:     "valid with HTTP 2.0",
			parts:    []string{"POST", "/api/users", "HTTP/2.0"},
			expected: true,
		},
		{
			name:     "less than three parts",
			parts:    []string{"GET", "/path"},
			expected: false,
		},
		{
			name:     "more than three parts",
			parts:    []string{"GET", "/path", "HTTP/1.1", "extra"},
			expected: false,
		},
		{
			name:     "empty parts",
			parts:    []string{},
			expected: false,
		},
		{
			name:     "one part",
			parts:    []string{"GET"},
			expected: false,
		},
		{
			name:     "three parts without HTTP prefix",
			parts:    []string{"GET", "/path", "1.1"},
			expected: false,
		},
		{
			name:     "three parts with lowercase http",
			parts:    []string{"GET", "/path", "http/1.1"},
			expected: false,
		},
		{
			name:     "three parts with partial HTTP",
			parts:    []string{"GET", "/path", "HTT/1.1"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validRequestLine(tc.parts)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestParseRequestLine(t *testing.T) {
	requestLineBytes := []byte("GET /users HTTP/1.1")

	result, err := ParseRequestLine(requestLineBytes)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("ParseRequestLine() returned nil result, expected non-nil request line")
	}
}

func TestParseRequestLineWithWhitespace(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{name: "leading spaces", input: "  GET /users HTTP/1.1"},
		{name: "trailing spaces", input: "GET /users HTTP/1.1  "},
		{name: "leading and trailing", input: "  GET /users HTTP/1.1  "},
		{name: "multiple spaces between parts", input: "GET    /users    HTTP/1.1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseRequestLine([]byte(tc.input))
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if result == nil {
				t.Error("ParseRequestLine() returned nil result, expected non-nil request line")
			}
		})
	}
}

func TestParseRequestLineNonASCII(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{name: "unicode in method", input: []byte("GÉT /users HTTP/1.1")},
		{name: "unicode in path", input: []byte("GET /üsers HTTP/1.1")},
		{name: "unicode in version", input: []byte("GET /users HTTÞ/1.1")},
		{name: "emoji", input: []byte("GET /users🎉 HTTP/1.1")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseRequestLine(tc.input)
			if err == nil {
				t.Error("expected error for non-ASCII input")
			}

			if result != nil {
				t.Error("ParseRequestLine() returned non-nil result on error, expected nil")
			}

			var errInvalidRequestLine pkgerrors.ErrInvalidRequestLine
			if !errors.As(err, &errInvalidRequestLine) {
				t.Errorf("expected ErrInvalidRequestLine, got %T", err)
			}
		})
	}
}

func TestParseRequestLineInvalidFormat(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{name: "two parts", input: "GET /users"},
		{name: "one part", input: "GET"},
		{name: "four parts", input: "GET /users HTTP/1.1 extra"},
		{name: "empty string", input: ""},
		{name: "only whitespace", input: "   "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseRequestLine([]byte(tc.input))
			if err == nil {
				t.Error("expected error for invalid format")
			}

			if result != nil {
				t.Error("ParseRequestLine() returned non-nil result on error, expected nil")
			}

			var errInvalidRequestLine pkgerrors.ErrInvalidRequestLine
			if !errors.As(err, &errInvalidRequestLine) {
				t.Errorf("expected ErrInvalidRequestLine, got %T", err)
			}
		})
	}
}

func TestParseRequestLineInvalidHTTPVersion(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{name: "no HTTP prefix", input: "GET /users 1.1"},
		{name: "lowercase http", input: "GET /users http/1.1"},
		{name: "partial HTTP", input: "GET /users HTT/1.1"},
		{name: "wrong prefix", input: "GET /users FTP/1.1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseRequestLine([]byte(tc.input))
			if err == nil {
				t.Error("expected error for invalid HTTP version")
			}

			if result != nil {
				t.Error("ParseRequestLine() returned non-nil result on error, expected nil")
			}

			var errInvalidRequestLine pkgerrors.ErrInvalidRequestLine
			if !errors.As(err, &errInvalidRequestLine) {
				t.Errorf("expected ErrInvalidRequestLine, got %T", err)
			}
		})
	}
}

func TestParseRequestLineDifferentMethods(t *testing.T) {
	testCases := []struct {
		name   string
		method string
	}{
		{name: "GET", method: "GET"},
		{name: "POST", method: "POST"},
		{name: "PUT", method: "PUT"},
		{name: "DELETE", method: "DELETE"},
		{name: "PATCH", method: "PATCH"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.method + " /users HTTP/1.1"
			result, err := ParseRequestLine([]byte(input))
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if result == nil {
				t.Error("ParseRequestLine() returned nil result, expected non-nil request line")
			}
		})
	}
}

func TestParseRequestLineDifferentHTTPVersions(t *testing.T) {
	testCases := []struct {
		name    string
		version string
	}{
		{name: "HTTP 1.0", version: "HTTP/1.0"},
		{name: "HTTP 1.1", version: "HTTP/1.1"},
		{name: "HTTP 2.0", version: "HTTP/2.0"},
		{name: "HTTP 3.0", version: "HTTP/3.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := "GET /users " + tc.version
			result, err := ParseRequestLine([]byte(input))
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if result == nil {
				t.Error("ParseRequestLine() returned nil result, expected non-nil request line")
			}
		})
	}
}
