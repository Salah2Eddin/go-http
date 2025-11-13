package uri

import (
	"testing"
)

func TestIsValidScheme(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid http", "http", true},
		{"valid https", "https", true},
		{"with digit", "http2", true},
		{"empty", "", false},
		{"starts with digit", "1http", false},
		{"contains invalid char", "htt@p", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidScheme(tc.input)
			if result != tc.expected {
				t.Errorf("isValidScheme(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidPort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid 80", "80", true},
		{"valid 8080", "8080", true},
		{"max port", "65535", true},
		{"min port", "1", true},
		{"empty", "", false},
		{"too high", "65536", false},
		{"non-digit", "80a0", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidPort(tc.input)
			if result != tc.expected {
				t.Errorf("isValidPort(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidIPv4(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid", "192.168.1.1", true},
		{"localhost", "127.0.0.1", true},
		{"broadcast", "255.255.255.255", true},
		{"too few octets", "192.168.1", false},
		{"octet > 255", "192.168.256.1", false},
		{"leading zero", "192.168.01.1", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidIPv4(tc.input)
			if result != tc.expected {
				t.Errorf("isValidIPv4(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidIPv6(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"compressed", "2001:db8::1", true},
		{"loopback", "::1", true},
		{"all zeros", "::", true},
		{"too many segments", "2001:db8:1:2:3:4:5:6:7", false},
		{"multiple ::", "2001::db8::1", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidIPv6(tc.input)
			if result != tc.expected {
				t.Errorf("isValidIPv6(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidDomainName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple", "example.com", true},
		{"subdomain", "sub.example.com", true},
		{"with hyphen", "my-domain.com", true},
		{"single label", "localhost", true},
		{"starts hyphen", "-example.com", false},
		{"ends hyphen", "example-.com", false},
		{"too long label", "a" + string(make([]byte, 64)) + ".com", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidDomainName(tc.input)
			if result != tc.expected {
				t.Errorf("isValidDomainName(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidUserInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple user", "user", true},
		{"user password", "user:pass", true},
		{"empty", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidUserInfo(tc.input)
			if result != tc.expected {
				t.Errorf("isValidUserInfo(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidHost(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"domain", "example.com", true},
		{"ipv4", "192.168.1.1", true},
		{"ipv6", "[::1]", true},
		{"empty", "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidHost(tc.input)
			if result != tc.expected {
				t.Errorf("isValidHost(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidAuthority(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"host only", "example.com", true},
		{"with port", "example.com:8080", true},
		{"with userinfo", "user@example.com", true},
		{"full", "user@example.com:8080", true},
		{"empty", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidAuthority(tc.input)
			if result != tc.expected {
				t.Errorf("isValidAuthority(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple", "/path", true},
		{"with query", "/path?key=value", true},
		{"with fragment", "/path#section", true},
		{"full", "/path?key=value#section", true},
		{"empty", "", true},
		{"fragment before query", "/path#section?key=value", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidPath(tc.input)
			if result != tc.expected {
				t.Errorf("isValidPath(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsAbsoluteURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"full uri", "http://example.com/path", true},
		{"with port", "http://example.com:8080/path", true},
		{"with query", "http://example.com/path?key=value", true},
		{"no scheme", "//example.com/path", false},
		{"relative", "/path", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isAbsoluteURI(tc.input)
			if result != tc.expected {
				t.Errorf("isAbsoluteURI(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsRelativeURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"absolute path", "/path", true},
		{"relative path", "path", true},
		{"with query", "/path?key=value", true},
		{"with fragment", "/path#frag", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isRelativeURI(tc.input)
			if result != tc.expected {
				t.Errorf("isRelativeURI(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestValidateURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"absolute", "http://example.com/path", true},
		{"relative", "/path", true},
		{"with query", "http://example.com?key=value", true},
		{"empty", "", false},
		{"invalid scheme", "ht!tp://example.com", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateURI(tc.input)
			if result != tc.expected {
				t.Errorf("ValidateURI(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
