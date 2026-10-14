package router

import (
	"testing"
)

func TestIsWildcard(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"wildcard", "*", true},
		{"empty", "", false},
		{"single letter", "a", false},
		{"wildcard with letter", "*a", false},
		{"letter with wildcard", "a*", false},
		{"double wildcard", "**", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isWildcard(tc.input)
			if result != tc.expected {
				t.Errorf("isWildcard(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
