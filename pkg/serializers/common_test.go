package serializers

import (
	"testing"
)

func TestNeedQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty", "", false},
		{"alphanum", "validtoken123", false},
		{"space", "invalid token", true},
		{"comma", "invalid,token", true},
		{"at sign", "invalid@token", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := needQuotes(tc.input)
			if result != tc.expected {
				t.Errorf("needQuotes(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestQuoteString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", `""`},
		{"simple", "hello", `"hello"`},
		{"single backslash", "a\\b", `"a\\b"`},
		{"single quote", "a\"b", `"a\"b"`},
		{"backslash and quote", "a\\\"b", `"a\\\"b"`},
		{"multiple backslashes", "\\\\", `"\\\\"`},
		{"multiple quotes", "\"\"", `"\"\""`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := quoteString(tc.input)
			if result != tc.expected {
				t.Errorf("quoteString(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}
