package httpheaders

import (
	"testing"
)

func TestValidHeaderName(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected bool
	}{
		// Valid cases
		{
			name:     "simple header name",
			input:    []byte("Content-Type"),
			expected: true,
		},
		{
			name:     "single character",
			input:    []byte("X"),
			expected: true,
		},
		{
			name:     "with numbers",
			input:    []byte("X-Custom-123"),
			expected: true,
		},
		{
			name:     "with special characters",
			input:    []byte("X-Custom!#$_"),
			expected: true,
		},

		// Invalid - whitespace
		{
			name:     "empty name",
			input:    []byte(""),
			expected: false,
		},
		{
			name:     "space in middle",
			input:    []byte("Content Type"),
			expected: false,
		},
		{
			name:     "space at start",
			input:    []byte(" Content-Type"),
			expected: false,
		},
		{
			name:     "tab character",
			input:    []byte("Content\tType"),
			expected: false,
		},
		{
			name:     "newline",
			input:    []byte("Content\nType"),
			expected: false,
		},
		{
			name:     "carriage return",
			input:    []byte("Content\rType"),
			expected: false,
		},

		// Invalid - non-ASCII and control characters
		{
			name:     "non-ASCII character",
			input:    []byte("Conten\xC3\xA9-Type"),
			expected: false,
		},
		{
			name:     "NUL character",
			input:    []byte("Content\x00Type"),
			expected: false,
		},
		{
			name:     "DEL character",
			input:    []byte("Content\x7FType"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validHeaderName(tc.input)
			if result != tc.expected {
				t.Errorf("validHeaderName(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestValidHeaderValue(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected bool
	}{
		// Valid cases
		{
			name:     "simple value",
			input:    []byte("application/json"),
			expected: true,
		},
		{
			name:     "value with parameters",
			input:    []byte("text/html; charset=utf-8"),
			expected: true,
		},
		{
			name:     "multiple values",
			input:    []byte("application/json, text/html"),
			expected: true,
		},
		{
			name:     "complex value",
			input:    []byte("text/html, application/json; q=0.9, */*; q=0.8"),
			expected: true,
		},
		{
			name:     "with spaces",
			input:    []byte("text html with spaces"),
			expected: true,
		},
		{
			name:     "with special characters",
			input:    []byte(`!@#$%^&*()_+-="quoted-string"`),
			expected: true,
		},
		{
			name:     "empty value",
			input:    []byte(""),
			expected: true,
		},

		// Invalid - CR, LF, NUL (CTL characters)
		{
			name:     "carriage return",
			input:    []byte("value\rvalue"),
			expected: false,
		},
		{
			name:     "line feed",
			input:    []byte("value\nvalue"),
			expected: false,
		},
		{
			name:     "NUL character",
			input:    []byte("value\x00value"),
			expected: false,
		},
		{
			name:     "BEL character (0x07)",
			input:    []byte("value\x07"),
			expected: false,
		},
		{
			name:     "BS character (0x08)",
			input:    []byte("value\x08"),
			expected: false,
		},
		{
			name:     "TAB character (0x09)",
			input:    []byte("value\x09"),
			expected: false,
		},
		{
			name:     "VT character (0x0B)",
			input:    []byte("value\x0B"),
			expected: false,
		},
		{
			name:     "FF character (0x0C)",
			input:    []byte("value\x0C"),
			expected: false,
		},
		{
			name:     "ESC character (0x1B)",
			input:    []byte("value\x1B"),
			expected: false,
		},
		{
			name:     "DEL character (0x7F)",
			input:    []byte("value\x7F"),
			expected: false,
		},
		{
			name:     "multiple CTL characters",
			input:    []byte("value\x01\x02\x03"),
			expected: false,
		},
		{
			name:     "CR+LF sequence",
			input:    []byte("value\r\nvalue"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validHeaderValue(tc.input)
			if result != tc.expected {
				t.Errorf("validHeaderValue(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}
