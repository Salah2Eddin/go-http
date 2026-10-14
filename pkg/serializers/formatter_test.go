package serializers

import (
	"testing"
)

func TestFormatterFactory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"content-type", "content-type", "identity"},
		{"default", "default", "quoted"},
		{"unknown", "unknown-key", "quoted"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := formatterFactory(tc.input)
			switch tc.expected {
			case "identity":
				if _, ok := result.(IdentityFormatter); !ok {
					t.Errorf("formatterFactory(%q) returned %T, expected IdentityFormatter", tc.input, result)
				}
			case "quoted":
				if _, ok := result.(QuotedFormatter); !ok {
					t.Errorf("formatterFactory(%q) returned %T, expected QuotedFormatter", tc.input, result)
				}
			}
		})
	}
}

func TestQuotedFormatterFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"needs quotes with space", "hello world", `"hello world"`},
		{"no quotes needed", "validtoken", "validtoken"},
		{"empty", "", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatter := QuotedFormatter{}
			result := formatter.format(tc.input)
			if result != tc.expected {
				t.Errorf("QuotedFormatter.format(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIdentityFormatterFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "hello", "hello"},
		{"with space", "hello world", "hello world"},
		{"with special chars", "hello@world", "hello@world"},
		{"empty", "", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatter := IdentityFormatter{}
			result := formatter.format(tc.input)
			if result != tc.expected {
				t.Errorf("IdentityFormatter.format(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}
