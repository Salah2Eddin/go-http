package charutil

import (
	"testing"
)

func TestIsUnreserved(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		// Unreserved characters: ALPHA / DIGIT / "-" / "." / "_" / "~"
		{"lowercase a", 'a', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', true},
		{"digit 9", '9', true},
		{"hyphen", '-', true},
		{"dot", '.', true},
		{"underscore", '_', true},
		{"tilde", '~', true},
		// Not unreserved
		{"space", ' ', false},
		{"exclamation", '!', false},
		{"at sign", '@', false},
		{"colon", ':', false},
		{"slash", '/', false},
		{"question mark", '?', false},
		{"hash", '#', false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsUnreserved(tc.input)
			if result != tc.expected {
				t.Errorf("IsUnreserved(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsSubDelim(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		// Sub-delimiters: ! $ & ' ( ) * + , ; =
		{"exclamation", '!', true},
		{"dollar", '$', true},
		{"ampersand", '&', true},
		{"single quote", '\'', true},
		{"left paren", '(', true},
		{"right paren", ')', true},
		{"asterisk", '*', true},
		{"plus", '+', true},
		{"comma", ',', true},
		{"semicolon", ';', true},
		{"equals", '=', true},
		// Not sub-delim
		{"lowercase a", 'a', false},
		{"uppercase Z", 'Z', false},
		{"digit 0", '0', false},
		{"hyphen", '-', false},
		{"dot", '.', false},
		{"underscore", '_', false},
		{"at sign", '@', false},
		{"colon", ':', false},
		{"space", ' ', false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsSubDelim(tc.input)
			if result != tc.expected {
				t.Errorf("IsSubDelim(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsPChar(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		// PChar = unreserved | sub-delim | ":" | "@"
		// Unreserved
		{"lowercase a", 'a', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', true},
		{"hyphen", '-', true},
		{"dot", '.', true},
		{"underscore", '_', true},
		{"tilde", '~', true},
		// Sub-delim
		{"exclamation", '!', true},
		{"dollar", '$', true},
		{"ampersand", '&', true},
		{"equals", '=', true},
		// Special pchar characters
		{"colon", ':', true},
		{"at sign", '@', true},
		// Not pchar
		{"space", ' ', false},
		{"slash", '/', false},
		{"question mark", '?', false},
		{"hash", '#', false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsPChar(tc.input)
			if result != tc.expected {
				t.Errorf("IsPChar(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsSchemeChar(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		// Scheme characters: ALPHA / DIGIT / "+" / "-" / "."
		{"lowercase a", 'a', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', true},
		{"digit 9", '9', true},
		{"plus", '+', true},
		{"hyphen", '-', true},
		{"dot", '.', true},
		// Not scheme char
		{"underscore", '_', false},
		{"colon", ':', false},
		{"at sign", '@', false},
		{"space", ' ', false},
		{"exclamation", '!', false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsSchemeChar(tc.input)
			if result != tc.expected {
				t.Errorf("IsSchemeChar(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsUserInfoChar(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		// UserInfo characters: unreserved | sub-delim | ":" | "%"
		// Unreserved
		{"lowercase a", 'a', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', true},
		{"hyphen", '-', true},
		{"dot", '.', true},
		{"underscore", '_', true},
		{"tilde", '~', true},
		// Sub-delim
		{"exclamation", '!', true},
		{"dollar", '$', true},
		{"ampersand", '&', true},
		{"equals", '=', true},
		// Special userinfo characters
		{"colon", ':', true},
		{"percent", '%', true},
		// Not userinfo char
		{"space", ' ', false},
		{"at sign", '@', false},
		{"slash", '/', false},
		{"null byte", 0x00, false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsUserInfoChar(tc.input)
			if result != tc.expected {
				t.Errorf("IsUserInfoChar(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
