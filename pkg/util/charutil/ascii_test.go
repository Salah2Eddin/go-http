package charutil

import (
	"testing"
)

func TestIsVisibleASCII(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"space (not visible)", AsciiSpace, false},
		{"tab (not visible)", AsciiTab, false},
		{"exclamation mark", '!', true},
		{"lowercase a", 'a', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', true},
		{"digit 9", '9', true},
		{"tilde", '~', true},
		{"DEL (not visible)", AsciiDelete, false},
		{"null byte", 0x00, false},
		{"below visible range", 0x20, false},
		{"above visible range", 0x7F, false},
		{"non-ASCII", 0x80, false},
		{"non-ASCII high", 0xFF, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsVisibleASCII(tc.input)
			if result != tc.expected {
				t.Errorf("IsVisibleASCII(0x%02X) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsASCII(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"null byte", 0x00, true},
		{"space", AsciiSpace, true},
		{"tab", AsciiTab, true},
		{"lowercase a", 'a', true},
		{"uppercase Z", 'Z', true},
		{"digit 0", '0', true},
		{"digit 9", '9', true},
		{"DEL", AsciiDelete, true},
		{"max ASCII", AsciiMax, true},
		{"min ASCII", AsciiMin, true},
		{"non-ASCII low", 0x80, false},
		{"non-ASCII mid", 0xC0, false},
		{"non-ASCII high", 0xFF, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsASCII(tc.input)
			if result != tc.expected {
				t.Errorf("IsASCII(0x%02X) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsCTLCharASCII(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"null byte", 0x00, true},
		{"BEL", 0x07, true},
		{"TAB", AsciiTab, true},
		{"LF", 0x0A, true},
		{"CR", 0x0D, true},
		{"max control", AsciiCtlMax, true},
		{"space (not control)", AsciiSpace, false},
		{"exclamation mark", '!', false},
		{"lowercase a", 'a', false},
		{"DEL", AsciiDelete, true},
		{"non-ASCII", 0x80, false},
		{"above DEL", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsCTLCharASCII(tc.input)
			if result != tc.expected {
				t.Errorf("IsCTLCharASCII(0x%02X) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsWhiteSpaceASCII(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"space", AsciiSpace, true},
		{"tab", AsciiTab, true},
		{"null byte", 0x00, false},
		{"LF", 0x0A, false},
		{"CR", 0x0D, false},
		{"exclamation mark", '!', false},
		{"lowercase a", 'a', false},
		{"digit 0", '0', false},
		{"non-ASCII", 0x80, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsWhiteSpaceASCII(tc.input)
			if result != tc.expected {
				t.Errorf("IsWhiteSpaceASCII(0x%02X) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestValidateAsciiEncoding(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{"empty slice", []byte{}, true},
		{"single ASCII", []byte{'a'}, true},
		{"multiple ASCII", []byte{'H', 'e', 'l', 'l', 'o'}, true},
		{"ASCII with spaces", []byte{'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd'}, true},
		{"all ASCII range", []byte{0x00, 0x7F}, true},
		{"single non-ASCII", []byte{0x80}, false},
		{"non-ASCII at start", []byte{0x80, 'a', 'b'}, false},
		{"non-ASCII in middle", []byte{'a', 0x80, 'b'}, false},
		{"non-ASCII at end", []byte{'a', 'b', 0x80}, false},
		{"multiple non-ASCII", []byte{0x80, 0xFF, 0xC0}, false},
		{"mixed ASCII and non-ASCII", []byte{'a', 0x80, 'b', 0xFF}, false},
		{"UTF-8 start byte", []byte{0xC0}, false},
		{"UTF-8 continuation", []byte{0x80}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateAsciiEncoding(tc.input)
			if result != tc.expected {
				t.Errorf("ValidateAsciiEncoding(%v) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
