package httpheaders

import (
	"testing"
)

func TestNewHeaders(t *testing.T) {
	headers := New()
	if headers == nil {
		t.Error("New() returned nil, expected non-nil headers instance")
	}
	if headers.headers == nil {
		t.Error("headers.headers map is nil, expected non-nil map")
	}
	if len(headers.headers) != 0 {
		t.Errorf("expected empty headers map, got %d headers", len(headers.headers))
	}
}

func TestAddFromStringSimple(t *testing.T) {
	headers := New()
	err := headers.AddFromString("Content-Type", "application/json")

	if err != nil {
		t.Errorf("AddFromString() error = %v", err)
	}

	retrieved, exists := headers.Get("Content-Type")
	if !exists {
		t.Error("Get('Content-Type') returned exists=false, expected header to exist after AddFromString")
	}
	if retrieved == nil {
		t.Error("Get() returned nil header, expected non-nil header")
	}

	values := retrieved.Values()
	if len(values) != 1 {
		t.Errorf("expected exactly 1 value, got %d", len(values))
	}
}

func TestAddFromStringWithParameters(t *testing.T) {
	headers := New()
	err := headers.AddFromString("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		t.Errorf("AddFromString() with parameters error = %v", err)
	}

	retrieved, exists := headers.Get("Content-Type")
	if !exists {
		t.Error("Get('Content-Type') returned exists=false, expected header to exist after AddFromString")
	}

	values := retrieved.Values()
	if len(values) != 1 {
		t.Errorf("expected exactly 1 value with parameters, got %d", len(values))
	}
}

func TestAddFromStringMultipleValues(t *testing.T) {
	headers := New()
	err := headers.AddFromString("Accept", "application/json, text/html, text/plain")

	if err != nil {
		t.Errorf("AddFromString() with multiple values error = %v", err)
	}

	retrieved, exists := headers.Get("Accept")
	if !exists {
		t.Error("Get('Accept') returned exists=false, expected header to exist after AddFromString")
	}

	values := retrieved.Values()
	if len(values) != 3 {
		t.Errorf("expected exactly 3 values, got %d", len(values))
	}
}

func TestAddFromStringMultipleValuesWithParameters(t *testing.T) {
	headers := New()
	err := headers.AddFromString("Accept", "application/json; q=0.9, text/html; q=0.8")

	if err != nil {
		t.Errorf("AddFromString() with multiple values and parameters error = %v", err)
	}

	retrieved, exists := headers.Get("Accept")
	if !exists {
		t.Error("Get('Accept') returned exists=false, expected header to exist after AddFromString")
	}

	values := retrieved.Values()
	if len(values) != 2 {
		t.Errorf("expected exactly 2 values, got %d", len(values))
	}
}

func TestAddFromHeader(t *testing.T) {
	headers := New()

	header, _ := NewHeaderFromString("Content-Type", "application/json")
	headers.AddFromHeader(header)

	retrieved, exists := headers.Get("Content-Type")
	if !exists {
		t.Error("Get('Content-Type') returned exists=false, expected header to exist after AddFromString")
	}
	if retrieved == nil {
		t.Error("Get() returned nil header, expected non-nil header")
	}
}

func TestAddFromHeaderCaseInsensitive(t *testing.T) {
	headers := New()

	header, _ := NewHeaderFromString("Content-Type", "application/json")
	headers.AddFromHeader(header)

	// Test retrieval with different cases
	cases := []string{"Content-Type", "content-type", "CONTENT-TYPE", "CoNtEnT-tYpE"}
	for _, c := range cases {
		_, exists := headers.Get(c)
		if !exists {
			t.Errorf("Get(%q) returned exists=false, expected to find header with case-insensitive lookup", c)
		}
	}
}

func TestAddFromHeaderMergesValues(t *testing.T) {
	headers := New()

	header1, _ := NewHeaderFromString("Accept", "application/json")
	headers.AddFromHeader(header1)

	header2, _ := NewHeaderFromString("Accept", "text/html")
	headers.AddFromHeader(header2)

	retrieved, exists := headers.Get("Accept")
	if !exists {
		t.Error("Get('Accept') returned exists=false, expected header to exist after AddFromString")
	}

	values := retrieved.Values()
	if len(values) != 2 {
		t.Errorf("expected exactly 2 values after merge, got %d", len(values))
	}
}

func TestAddFromHeaderMergesMultipleValues(t *testing.T) {
	headers := New()

	header1, _ := NewHeaderFromString("Accept", "application/json, text/html")
	headers.AddFromHeader(header1)

	header2, _ := NewHeaderFromString("Accept", "text/plain, image/png")
	headers.AddFromHeader(header2)

	retrieved, exists := headers.Get("Accept")
	if !exists {
		t.Error("Get('Accept') returned exists=false, expected header to exist after AddFromString")
	}

	values := retrieved.Values()
	if len(values) != 4 {
		t.Errorf("expected exactly 4 total values after merge, got %d", len(values))
	}
}

func TestAddFromHeaderMergesValuesWithParameters(t *testing.T) {
	headers := New()

	header1, _ := NewHeaderFromString("Accept", "application/json; q=0.9")
	headers.AddFromHeader(header1)

	header2, _ := NewHeaderFromString("Accept", "text/html; q=0.8")
	headers.AddFromHeader(header2)

	retrieved, exists := headers.Get("Accept")
	if !exists {
		t.Error("Get('Accept') returned exists=false, expected header to exist after AddFromString")
	}

	values := retrieved.Values()
	if len(values) != 2 {
		t.Errorf("expected exactly 2 values, got %d", len(values))
	}
}

func TestGet(t *testing.T) {
	headers := New()

	// Test getting non-existent header
	_, exists := headers.Get("Non-Existent")
	if exists {
		t.Error("Get('Non-Existent') returned exists=true, expected header to not exist")
	}

	// Add a header and retrieve it
	headers.AddFromString("X-Custom-Header", "value1")
	retrieved, exists := headers.Get("X-Custom-Header")
	if !exists {
		t.Error("Get('X-Custom-Header') returned exists=false, expected header to exist after AddFromString")
	}
	if retrieved == nil {
		t.Error("Get() returned nil header, expected non-nil header")
	}
}

func TestGetCaseInsensitive(t *testing.T) {
	headers := New()
	headers.AddFromString("X-Custom-Header", "value1")

	testCases := []string{"X-Custom-Header", "x-custom-header", "X-CUSTOM-HEADER"}
	for _, tc := range testCases {
		_, exists := headers.Get(tc)
		if !exists {
			t.Errorf("expected to find header with key %q", tc)
		}
	}
}

func TestHeadersMethod(t *testing.T) {
	headers := New()

	// Empty headers
	result := headers.Headers()
	if len(result) != 0 {
		t.Errorf("Headers() returned %d headers, expected empty map", len(result))
	}

	// Add some headers
	headers.AddFromString("Content-Type", "application/json")
	headers.AddFromString("Authorization", "Bearer token")

	result = headers.Headers()
	if len(result) != 2 {
		t.Errorf("expected 2 headers, got %d", len(result))
	}

	// Check that both headers are present
	if _, exists := result["content-type"]; !exists {
		t.Error("Headers() map missing 'content-type' key, expected it to be present")
	}
	if _, exists := result["authorization"]; !exists {
		t.Error("Headers() map missing 'authorization' key, expected it to be present")
	}
}

func TestMultipleHeaderOperations(t *testing.T) {
	headers := New()

	// Add multiple different headers
	headers.AddFromString("Content-Type", "application/json; charset=utf-8")
	headers.AddFromString("Accept", "application/json, text/html")
	headers.AddFromString("User-Agent", "TestAgent/1.0")

	// Verify all exist
	if _, exists := headers.Get("Content-Type"); !exists {
		t.Error("Get('Content-Type') returned exists=false, expected header to exist")
	}
	if _, exists := headers.Get("Accept"); !exists {
		t.Error("Get('Accept') returned exists=false, expected header to exist")
	}
	if _, exists := headers.Get("User-Agent"); !exists {
		t.Error("Get('User-Agent') returned exists=false, expected header to exist")
	}

	// Verify total count
	if len(headers.Headers()) != 3 {
		t.Errorf("expected 3 headers, got %d", len(headers.Headers()))
	}
}

func TestAddHeaderWithEmptyValue(t *testing.T) {
	headers := New()
	err := headers.AddFromString("X-Empty", "")

	if err == nil {
		t.Error("AddFromString() with empty value returned nil error, expected non-nil error")
	}

	_, exists := headers.Get("X-Empty")
	if exists {
		t.Error("Get('X-Empty') returned exists=true, expected header with empty value to not exist")
	}
}

func TestAddFromStringComplexScenario(t *testing.T) {
	headers := New()

	// Simulate HTTP Accept header with quality values
	err := headers.AddFromString("Accept", "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8")
	if err != nil {
		t.Errorf("AddFromString() with complex Accept header error = %v", err)
	}

	retrieved, exists := headers.Get("Accept")
	if !exists {
		t.Error("Get('Accept') returned exists=false, expected header to exist after AddFromString")
	}

	values := retrieved.Values()
	if len(values) != 4 {
		t.Errorf("expected exactly 4 values in complex Accept header, got %d", len(values))
	}
}

func TestAddFromStringContentTypeVariants(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{
			name:  "simple content type",
			value: "text/plain",
		},
		{
			name:  "with charset",
			value: "text/html; charset=utf-8",
		},
		{
			name:  "with boundary",
			value: "multipart/form-data; boundary=----WebKitFormBoundary",
		},
		{
			name:  "with multiple parameters",
			value: "text/plain; charset=utf-8; format=flowed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := New()
			err := headers.AddFromString("Content-Type", tc.value)

			if err != nil {
				t.Errorf("AddFromString() error = %v", err)
			}

			retrieved, exists := headers.Get("Content-Type")
			if !exists {
				t.Error("Get('Content-Type') returned exists=false, expected header to exist after AddFromString")
			}

			values := retrieved.Values()
			if len(values) != 1 {
				t.Errorf("expected exactly 1 value, got %d", len(values))
			}
		})
	}
}
