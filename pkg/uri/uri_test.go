package uri

import (
	"testing"
)

func TestNewUri(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple http", "http://example.com", "http://example.com"},
		{"with path", "http://example.com/path", "http://example.com/path"},
		{"with port", "http://example.com:8080/path", "http://example.com:8080/path"},
		{"with query", "http://example.com/path?key=value", "http://example.com/path?key=value"},
		{"with fragment", "http://example.com/path#section", "http://example.com/path#section"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NewUri(tc.input)
			if result == nil {
				t.Errorf("NewUri(%q) returned nil", tc.input)
			}
		})
	}
}

func TestUriString(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "http://example.com"},
		{"with path", "http://example.com/path"},
		{"with port", "http://example.com:8080"},
		{"full uri", "http://user@example.com:8080/path?key=value#frag"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uri := NewUri(tc.input)
			result := uri.String()
			if result == "" {
				t.Errorf("String() returned empty")
			}
		})
	}
}

func TestGetQueryParameter(t *testing.T) {
	tests := []struct {
		name          string
		uri           string
		param         string
		expectedFound bool
	}{
		{"existing param", "http://example.com?key=value", "key", true},
		{"missing param", "http://example.com?key=value", "other", false},
		{"no query", "http://example.com", "key", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uri := NewUri(tc.uri)
			_, found := uri.GetQueryParameter(tc.param)
			if found != tc.expectedFound {
				t.Errorf("GetQueryParameter(%q) found=%v, want %v", tc.param, found, tc.expectedFound)
			}
		})
	}
}

func TestGetSegments(t *testing.T) {
	tests := []struct {
		name          string
		uri           string
		expectedCount int
	}{
		{"root path", "http://example.com/", 2},
		{"single segment", "http://example.com/path", 2},
		{"multiple segments", "http://example.com/path/to/resource", 4},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uri := NewUri(tc.uri)
			segments := uri.GetSegments()
			if len(segments) != tc.expectedCount {
				t.Errorf("GetSegments() returned %d segments, want %d", len(segments), tc.expectedCount)
			}
		})
	}
}
