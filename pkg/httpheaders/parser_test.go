package httpheaders

import (
	"bytes"
	"testing"
)

// processHeaderName tests
func TestProcessHeaderName(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "simple name",
			input:    []byte("Content-Type"),
			expected: "content-type",
		},
		{
			name:     "all uppercase",
			input:    []byte("AUTHORIZATION"),
			expected: "authorization",
		},
		{
			name:     "mixed case",
			input:    []byte("X-Custom-Header"),
			expected: "x-custom-header",
		},
		{
			name:     "already lowercase",
			input:    []byte("accept"),
			expected: "accept",
		},
		{
			name:     "single character",
			input:    []byte("X"),
			expected: "x",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := processHeaderName(tc.input)
			if result != tc.expected {
				t.Errorf("processHeaderName(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

// nameValueSplit tests
func TestNameValueSplit(t *testing.T) {
	testCases := []struct {
		name          string
		input         []byte
		expectedName  []byte
		expectedValue []byte
		expectedFound bool
	}{
		{
			name:          "simple header",
			input:         []byte("Content-Type: application/json"),
			expectedName:  []byte("Content-Type"),
			expectedValue: []byte(" application/json"),
			expectedFound: true,
		},
		{
			name:          "no space after colon",
			input:         []byte("Authorization:Bearer token"),
			expectedName:  []byte("Authorization"),
			expectedValue: []byte("Bearer token"),
			expectedFound: true,
		},
		{
			name:          "empty value",
			input:         []byte("X-Header:"),
			expectedName:  []byte("X-Header"),
			expectedValue: []byte(""),
			expectedFound: true,
		},
		{
			name:          "no colon",
			input:         []byte("Invalid-Header"),
			expectedName:  nil,
			expectedValue: nil,
			expectedFound: false,
		},
		{
			name:          "multiple colons",
			input:         []byte("Key: value:with:colons"),
			expectedName:  []byte("Key"),
			expectedValue: []byte(" value:with:colons"),
			expectedFound: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name, value, found := nameValueSplit(tc.input)
			if found != tc.expectedFound {
				t.Errorf("nameValueSplit found = %v, expected %v", found, tc.expectedFound)
			}
			if found {
				if !bytes.Equal(name, tc.expectedName) {
					t.Errorf("nameValueSplit name = %q, expected %q", name, tc.expectedName)
				}
				if !bytes.Equal(value, tc.expectedValue) {
					t.Errorf("nameValueSplit value = %q, expected %q", value, tc.expectedValue)
				}
			}
		})
	}
}

// parseUnquotedValue tests
func TestParseUnquotedValue(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expected      []byte
		expectedError bool
	}{
		{
			name:          "simple value",
			input:         "application/json",
			expected:      []byte("application/json"),
			expectedError: false,
		},
		{
			name:          "value with leading spaces",
			input:         "   application/json",
			expected:      []byte("application/json"),
			expectedError: false,
		},
		{
			name:          "value with trailing spaces",
			input:         "application/json   ",
			expected:      []byte("application/json"),
			expectedError: false,
		},
		{
			name:          "value stops at comma",
			input:         "text/html, application/json",
			expected:      []byte("text/html"),
			expectedError: false,
		},
		{
			name:          "value stops at semicolon",
			input:         "text/html; charset=utf-8",
			expected:      []byte("text/html"),
			expectedError: false,
		},
		{
			name:          "internal spaces preserved",
			input:         "value with spaces,next",
			expected:      []byte("value with spaces"),
			expectedError: false,
		},
		{
			name:          "empty input",
			input:         "",
			expected:      []byte(nil),
			expectedError: false,
		},
		{
			name:          "contains unquoted quote",
			input:         `value"notallowed`,
			expected:      nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tc.input))
			result, err := parseUnquotedValue(reader)

			if (err != nil) != tc.expectedError {
				t.Errorf("parseUnquotedValue error = %v, expectedError = %v", err, tc.expectedError)
			}
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("parseUnquotedValue = %q, expected %q", result, tc.expected)
			}
		})
	}
}

// parseQuotedValue tests
func TestParseQuotedValue(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expected      []byte
		expectedError bool
	}{
		{
			name:          "simple quoted value",
			input:         `"quoted-value"`,
			expected:      []byte("quoted-value"),
			expectedError: false,
		},
		{
			name:          "quoted value with spaces",
			input:         `"quoted value with spaces"`,
			expected:      []byte("quoted value with spaces"),
			expectedError: false,
		},
		{
			name:          "quoted value with escaped quote",
			input:         `"value with \" escaped"`,
			expected:      []byte(`value with \" escaped`),
			expectedError: false,
		},
		{
			name:          "quoted value followed by comma",
			input:         `"value1",`,
			expected:      []byte("value1"),
			expectedError: false,
		},
		{
			name:          "quoted value followed by semicolon",
			input:         `"value";charset=utf-8`,
			expected:      []byte("value"),
			expectedError: false,
		},
		{
			name:          "quoted value followed by spaces then comma",
			input:         `"value"  ,`,
			expected:      []byte("value"),
			expectedError: false,
		},
		{
			name:          "unclosed quoted value",
			input:         `"unclosed`,
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "quoted value followed by unquoted",
			input:         `"value"unquoted`,
			expected:      nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tc.input))
			result, err := parseQuotedValue(reader)

			if (err != nil) != tc.expectedError {
				t.Errorf("parseQuotedValue error = %v, expectedError = %v", err, tc.expectedError)
			}
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("parseQuotedValue = %q, expected %q", result, tc.expected)
			}
		})
	}
}

// parseParameters tests
func TestParseParameters(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expected      []byte
		expectedError bool
	}{
		{
			name:          "simple parameter",
			input:         ";charset=utf-8",
			expected:      []byte("charset=utf-8"),
			expectedError: false,
		},
		{
			name:          "parameter with spaces before semicolon",
			input:         "  ;charset=utf-8",
			expected:      []byte("charset=utf-8"),
			expectedError: false,
		},
		{
			name:          "parameter stops at comma",
			input:         ";charset=utf-8,next",
			expected:      []byte("charset=utf-8"),
			expectedError: false,
		},
		{
			name:          "parameter stops at next semicolon",
			input:         ";charset=utf-8;next=param",
			expected:      []byte("charset=utf-8"),
			expectedError: false,
		},
		{
			name:          "empty input",
			input:         "",
			expected:      []byte(nil),
			expectedError: false,
		},
		{
			name:          "invalid syntax no semicolon",
			input:         "charset=utf-8",
			expected:      nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tc.input))
			result, err := parseParameters(reader)

			if (err != nil) != tc.expectedError {
				t.Errorf("parseParameters error = %v, expectedError = %v", err, tc.expectedError)
			}
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("parseParameters = %q, expected %q", result, tc.expected)
			}
		})
	}
}

// parseNextValue tests
func TestParseNextValue(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedValue  []byte
		expectedParams []byte
		expectedError  bool
	}{
		{
			name:           "unquoted value no params",
			input:          "application/json",
			expectedValue:  []byte("application/json"),
			expectedParams: []byte(nil),
			expectedError:  false,
		},
		{
			name:           "unquoted value with params",
			input:          "text/html;charset=utf-8",
			expectedValue:  []byte("text/html"),
			expectedParams: []byte("charset=utf-8"),
			expectedError:  false,
		},
		{
			name:           "quoted value no params",
			input:          `"quoted-value"`,
			expectedValue:  []byte("quoted-value"),
			expectedParams: []byte(nil),
			expectedError:  false,
		},
		{
			name:           "quoted value with params",
			input:          `"value";charset=utf-8`,
			expectedValue:  []byte("value"),
			expectedParams: []byte("charset=utf-8"),
			expectedError:  false,
		},
		{
			name:           "leading whitespace",
			input:          "   text/html",
			expectedValue:  []byte("text/html"),
			expectedParams: []byte(nil),
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tc.input))
			value, params, err := parseNextValue(reader)

			if (err != nil) != tc.expectedError {
				t.Errorf("parseNextValue error = %v, expectedError = %v", err, tc.expectedError)
			}
			if !bytes.Equal(value, tc.expectedValue) {
				t.Errorf("parseNextValue value = %q, expected %q", value, tc.expectedValue)
			}
			if !bytes.Equal(params, tc.expectedParams) {
				t.Errorf("parseNextValue params = %q, expected %q", params, tc.expectedParams)
			}
		})
	}
}

// splitHeaderValues tests
func TestSplitHeaderValues(t *testing.T) {
	testCases := []struct {
		name          string
		input         []byte
		expectedCount int
		expectedError bool
	}{
		{
			name:          "single value",
			input:         []byte("application/json"),
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:          "multiple values",
			input:         []byte("application/json, text/html, text/plain"),
			expectedCount: 3,
			expectedError: false,
		},
		{
			name:          "values with parameters",
			input:         []byte("text/html;charset=utf-8, application/json"),
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "quoted values",
			input:         []byte(`"value1", "value2"`),
			expectedCount: 2,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			values, params, err := splitHeaderValues(tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("splitHeaderValues error = %v, expectedError = %v", err, tc.expectedError)
			}
			if len(values) != tc.expectedCount {
				t.Errorf("splitHeaderValues count = %d, expected %d", len(values), tc.expectedCount)
			}
			if len(values) != len(params) {
				t.Errorf("splitHeaderValues values and params count mismatch")
			}
		})
	}
}

// processHeaderValues tests
func TestProcessHeaderValues(t *testing.T) {
	testCases := []struct {
		name          string
		input         []byte
		expectedCount int
		expectedError bool
	}{
		{
			name:          "single value",
			input:         []byte("application/json"),
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:          "multiple values",
			input:         []byte("application/json, text/html"),
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "value with parameters",
			input:         []byte("text/html;charset=utf-8"),
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:          "empty string",
			input:         []byte(""),
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := processHeaderValues(tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("processHeaderValues error = %v, expectedError = %v", err, tc.expectedError)
			}
			if len(result) != tc.expectedCount {
				t.Errorf("processHeaderValues count = %d, expected %d", len(result), tc.expectedCount)
			}
		})
	}
}

// parseHeaderLine tests
func TestParseHeaderLine(t *testing.T) {
	testCases := []struct {
		name          string
		input         []byte
		expectedError bool
	}{
		{
			name:          "valid simple header",
			input:         []byte("Content-Type: application/json"),
			expectedError: false,
		},
		{
			name:          "valid header with parameters",
			input:         []byte("Content-Type: text/html; charset=utf-8"),
			expectedError: false,
		},
		{
			name:          "valid header with multiple values",
			input:         []byte("Accept: application/json, text/html"),
			expectedError: false,
		},
		{
			name:          "no colon delimiter",
			input:         []byte("InvalidHeader"),
			expectedError: true,
		},
		{
			name:          "invalid header name with space",
			input:         []byte("Content Type: application/json"),
			expectedError: true,
		},
		{
			name:          "invalid header value with CTL",
			input:         []byte("X-Header: value\x00invalid"),
			expectedError: true,
		},
		{
			name:          "empty value",
			input:         []byte("X-Header: "),
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseHeaderLine(tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("parseHeaderLine error = %v, expectedError = %v", err, tc.expectedError)
			}
			if !tc.expectedError && result == nil {
				t.Error("parseHeaderLine returned nil header for valid input")
			}
		})
	}
}

// ParseRequestHeaders tests
func TestParseRequestHeaders(t *testing.T) {
	testCases := []struct {
		name          string
		input         [][]byte
		expectedCount int
		expectedError bool
	}{
		{
			name: "single header",
			input: [][]byte{
				[]byte("Content-Type: application/json"),
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "multiple headers",
			input: [][]byte{
				[]byte("Content-Type: application/json"),
				[]byte("Accept: text/html"),
				[]byte("Authorization: Bearer token"),
			},
			expectedCount: 3,
			expectedError: false,
		},
		{
			name: "headers with same name merge values",
			input: [][]byte{
				[]byte("Accept: application/json"),
				[]byte("Accept: text/html"),
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "invalid header line",
			input: [][]byte{
				[]byte("Content-Type: application/json"),
				[]byte("Invalid-Header"),
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseRequestHeaders(tc.input)

			if (err != nil) != tc.expectedError {
				t.Errorf("ParseRequestHeaders error = %v, expectedError = %v", err, tc.expectedError)
			}
			if !tc.expectedError && result == nil {
				t.Error("ParseRequestHeaders returned nil")
			}
			if !tc.expectedError {
				if len(result.Headers()) != tc.expectedCount {
					t.Errorf("ParseRequestHeaders count = %d, expected %d", len(result.Headers()), tc.expectedCount)
				}
			}
		})
	}
}
