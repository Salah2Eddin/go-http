package uri

import (
	"testing"
)

func TestParseURI(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedScheme   string
		expectedHost     string
		expectedPort     string
		expectedPath     string
		expectedFragment string
		expectedAbsolute bool
	}{
		{
			"absolute with all components",
			"http://example.com:8080/path?key=value#frag",
			"http", "example.com", "8080", "/path", "frag", true,
		},
		{
			"absolute minimal",
			"http://example.com",
			"http", "example.com", "", "", "", true,
		},
		{
			"relative path",
			"/path",
			"", "", "", "/path", "", false,
		},
		{
			"with userinfo",
			"http://user@example.com/path",
			"http", "example.com", "", "/path", "", true,
		},
		{
			"ipv6",
			"http://[::1]:8080/path",
			"http", "[::1]", "8080", "/path", "", true,
		},
		{
			"empty string",
			"",
			"", "", "", "", "", false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := &Uri{}
			parseURI(tc.input, u)
			if u.scheme != tc.expectedScheme {
				t.Errorf("scheme = %q, want %q", u.scheme, tc.expectedScheme)
			}
			if u.host != tc.expectedHost {
				t.Errorf("host = %q, want %q", u.host, tc.expectedHost)
			}
			if u.port != tc.expectedPort {
				t.Errorf("port = %q, want %q", u.port, tc.expectedPort)
			}
			if u.path != tc.expectedPath {
				t.Errorf("path = %q, want %q", u.path, tc.expectedPath)
			}
			if u.fragment != tc.expectedFragment {
				t.Errorf("fragment = %q, want %q", u.fragment, tc.expectedFragment)
			}
			if u.isAbsolute != tc.expectedAbsolute {
				t.Errorf("isAbsolute = %v, want %v", u.isAbsolute, tc.expectedAbsolute)
			}
		})
	}
}

func TestParseAuthority(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedUser string
		expectedHost string
		expectedPort string
	}{
		{"host only", "example.com", "", "example.com", ""},
		{"with port", "example.com:8080", "", "example.com", "8080"},
		{"with userinfo", "user@example.com", "user", "example.com", ""},
		{"full", "user@example.com:8080", "user", "example.com", "8080"},
		{"ipv6 with port", "[::1]:8080", "", "[::1]", "8080"},
		{"ipv6 only", "[::1]", "", "[::1]", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := &Uri{}
			parseAuthority(tc.input, u)
			if u.userInfo != tc.expectedUser {
				t.Errorf("userInfo = %q, want %q", u.userInfo, tc.expectedUser)
			}
			if u.host != tc.expectedHost {
				t.Errorf("host = %q, want %q", u.host, tc.expectedHost)
			}
			if u.port != tc.expectedPort {
				t.Errorf("port = %q, want %q", u.port, tc.expectedPort)
			}
		})
	}
}

func TestExtractPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"path only", "/path", "/path"},
		{"with query", "/path?key=value", "/path"},
		{"with fragment", "/path#frag", "/path"},
		{"with both", "/path?query#frag", "/path"},
		{"empty", "", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractPath(tc.input)
			if result != tc.expected {
				t.Errorf("extractPath(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestExtractQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with query", "/path?key=value", "key=value"},
		{"with fragment after query", "/path?key=value#frag", "key=value"},
		{"no query", "/path", ""},
		{"empty query", "/path?", ""},
		{"fragment only", "/path#frag", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractQuery(tc.input)
			if result != tc.expected {
				t.Errorf("extractQuery(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestExtractFragment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with fragment", "/path#frag", "frag"},
		{"no fragment", "/path", ""},
		{"fragment with query", "/path?key#frag", "frag"},
		{"empty fragment", "/path#", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractFragment(tc.input)
			if result != tc.expected {
				t.Errorf("extractFragment(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedParams map[string]string
	}{
		{"single param", "key=value", map[string]string{"key": "value"}},
		{"multiple params", "key1=value1&key2=value2", map[string]string{"key1": "value1", "key2": "value2"}},
		{"empty", "", map[string]string{}},
		{"no value", "key=", map[string]string{"key": ""}},
		{"param without equals", "key", map[string]string{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := &Uri{query: make(map[string]string)}
			parseQuery(tc.input, u)
			if len(u.query) != len(tc.expectedParams) {
				t.Errorf("parseQuery returned %d params, want %d", len(u.query), len(tc.expectedParams))
			}
			for key, expectedVal := range tc.expectedParams {
				val, ok := u.query[key]
				if !ok {
					t.Errorf("parseQuery missing key %q", key)
				}
				if val != expectedVal {
					t.Errorf("query[%q] = %q, want %q", key, val, expectedVal)
				}
			}
		})
	}
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedPath     string
		expectedFragment string
		expectedQuery    bool
	}{
		{"/path?key=value#frag", "/path?key=value#frag", "/path", "frag", true},
		{"/path", "/path", "/path", "", false},
		{"?key=value", "?key=value", "", "", true},
		{"#frag", "#frag", "", "frag", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := &Uri{query: make(map[string]string)}
			parsePath(tc.input, u)
			if u.path != tc.expectedPath {
				t.Errorf("path = %q, want %q", u.path, tc.expectedPath)
			}
			if u.fragment != tc.expectedFragment {
				t.Errorf("fragment = %q, want %q", u.fragment, tc.expectedFragment)
			}
		})
	}
}
