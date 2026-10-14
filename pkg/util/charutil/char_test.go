package charutil

import (
	"testing"
)

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"lowercase a", 'a', true},
		{"lowercase z", 'z', true},
		{"lowercase m", 'm', true},
		{"uppercase A", 'A', true},
		{"uppercase Z", 'Z', true},
		{"uppercase M", 'M', true},
		{"digit 0", '0', false},
		{"digit 9", '9', false},
		{"space", ' ', false},
		{"exclamation mark", '!', false},
		{"below lowercase", 'a' - 1, false},
		{"above lowercase", 'z' + 1, false},
		{"below uppercase", 'A' - 1, false},
		{"above uppercase", 'Z' + 1, false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsAlpha(tc.input)
			if result != tc.expected {
				t.Errorf("IsAlpha(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"digit 0", '0', true},
		{"digit 5", '5', true},
		{"digit 9", '9', true},
		{"below 0", '0' - 1, false},
		{"above 9", '9' + 1, false},
		{"lowercase a", 'a', false},
		{"uppercase A", 'A', false},
		{"space", ' ', false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsDigit(tc.input)
			if result != tc.expected {
				t.Errorf("IsDigit(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsAlphaNum(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"lowercase a", 'a', true},
		{"lowercase z", 'z', true},
		{"uppercase A", 'A', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', true},
		{"digit 9", '9', true},
		{"space", ' ', false},
		{"exclamation mark", '!', false},
		{"at sign", '@', false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsAlphaNum(tc.input)
			if result != tc.expected {
				t.Errorf("IsAlphaNum(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsHexAlpha(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"lowercase a", 'a', true},
		{"lowercase f", 'f', true},
		{"lowercase e", 'e', true},
		{"uppercase A", 'A', true},
		{"uppercase F", 'F', true},
		{"uppercase C", 'C', true},
		{"below lowercase a", 'a' - 1, false},
		{"above lowercase f", 'f' + 1, false},
		{"below uppercase A", 'A' - 1, false},
		{"above uppercase F", 'F' + 1, false},
		{"digit 0", '0', false},
		{"digit 9", '9', false},
		{"lowercase g", 'g', false},
		{"uppercase G", 'G', false},
		{"space", ' ', false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsHexAlpha(tc.input)
			if result != tc.expected {
				t.Errorf("IsHexAlpha(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsHexDigit(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"digit 0", '0', true},
		{"digit 5", '5', true},
		{"digit 9", '9', true},
		{"lowercase a", 'a', true},
		{"lowercase f", 'f', true},
		{"uppercase A", 'A', true},
		{"uppercase F", 'F', true},
		{"below 0", '0' - 1, false},
		{"above 9", '9' + 1, false},
		{"lowercase g", 'g', false},
		{"uppercase G", 'G', false},
		{"space", ' ', false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsHexDigit(tc.input)
			if result != tc.expected {
				t.Errorf("IsHexDigit(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsTChar(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		// Alphanumeric (should be true)
		{"lowercase a", 'a', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', true},
		{"digit 9", '9', true},
		// TChar symbols
		{"exclamation", '!', true},
		{"hash", '#', true},
		{"dollar", '$', true},
		{"percent", '%', true},
		{"ampersand", '&', true},
		{"single quote", '\'', true},
		{"asterisk", '*', true},
		{"plus", '+', true},
		{"minus", '-', true},
		{"dot", '.', true},
		{"caret", '^', true},
		{"underscore", '_', true},
		{"grave", '`', true},
		{"pipe", '|', true},
		{"tilde", '~', true},
		// Not TChar
		{"space", ' ', false},
		{"at sign", '@', false},
		{"equals", '=', false},
		{"comma", ',', false},
		{"semicolon", ';', false},
		{"left paren", '(', false},
		{"right paren", ')', false},
		{"colon", ':', false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsTChar(tc.input)
			if result != tc.expected {
				t.Errorf("IsTChar(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
