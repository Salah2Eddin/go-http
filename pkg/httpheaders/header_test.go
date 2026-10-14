package httpheaders

import (
	"testing"
)

func TestNewHeader(t *testing.T) {
	var values []*Value
	header := NewHeader("Content-Type", values)

	if header == nil {
		t.Fatal("NewHeader() returned nil, expected non-nil header")
	}
	if header.Name() != "Content-Type" {
		t.Errorf("expected name 'Content-Type', got '%s'", header.Name())
	}
	if len(header.Values()) != 0 {
		t.Errorf("expected no values, got %d", len(header.Values()))
	}
}

func TestNewHeaderWithValues(t *testing.T) {
	values := []*Value{
		{}, // Assuming Value type exists and can be instantiated
		{},
	}
	header := NewHeader("Accept", values)

	if header == nil {
		t.Fatal("NewHeader() returned nil, expected non-nil header")
	}
	if header.Name() != "Accept" {
		t.Errorf("expected name 'Accept', got '%s'", header.Name())
	}
	if len(header.Values()) != 2 {
		t.Errorf("expected 2 values, got %d", len(header.Values()))
	}
}

func TestNewHeaderFromStringSimple(t *testing.T) {
	header, err := NewHeaderFromString("Content-Type", "application/json")

	if err != nil {
		t.Errorf("NewHeaderFromString() error = %v", err)
	}
	if header == nil {
		t.Fatal("NewHeaderFromString() returned nil header, expected non-nil header")
	}
	if header.Name() != "Content-Type" {
		t.Errorf("expected name 'Content-Type', got '%s'", header.Name())
	}
	if len(header.Values()) != 1 {
		t.Errorf("expected one value, got %d", len(header.Values()))
	}
}

func TestNewHeaderFromStringWithParameters(t *testing.T) {
	header, err := NewHeaderFromString("Content-Type", "text/html; charset=utf-8")

	if err != nil {
		t.Errorf("NewHeaderFromString() with parameters error = %v", err)
	}
	if header == nil {
		t.Fatal("NewHeaderFromString() returned nil header, expected non-nil header")
	}
	if len(header.Values()) != 1 {
		t.Errorf("expected exactly 1 value with parameters, got %d", len(header.Values()))
	}
}

func TestNewHeaderFromStringMultipleValues(t *testing.T) {
	header, err := NewHeaderFromString("Accept", "application/json, text/html, text/plain")

	if err != nil {
		t.Errorf("NewHeaderFromString() with multiple values error = %v", err)
	}
	if header == nil {
		t.Fatal("NewHeaderFromString() returned nil header, expected non-nil header")
	}
	if len(header.Values()) != 3 {
		t.Errorf("expected exactly 3 values, got %d", len(header.Values()))
	}
}

func TestNewHeaderFromStringMultipleValuesWithParameters(t *testing.T) {
	header, err := NewHeaderFromString("Accept", "application/json; q=0.9, text/html; q=0.8, text/plain; q=0.7")

	if err != nil {
		t.Errorf("NewHeaderFromString() with multiple values and parameters error = %v", err)
	}
	if header == nil {
		t.Fatal("NewHeaderFromString() returned nil header, expected non-nil header")
	}
	if len(header.Values()) != 3 {
		t.Errorf("expected exactly 3 values, got %d", len(header.Values()))
	}
}

func TestNewHeaderFromStringEmptyValue(t *testing.T) {
	header, err := NewHeaderFromString("X-Empty", "")

	if err == nil {
		t.Errorf("NewHeaderFromString() with empty value error = %v", err)
	}
	if header != nil {
		t.Fatalf("NewHeaderFromString() returned non-nil header %v, expected nil for empty value", header)
	}
}

func TestName(t *testing.T) {
	header := NewHeader("Authorization", []*Value{})

	if header.Name() != "Authorization" {
		t.Errorf("expected 'Authorization', got '%s'", header.Name())
	}
}

func TestValues(t *testing.T) {
	values := []*Value{{}, {}, {}}
	header := NewHeader("Accept", values)

	retrievedValues := header.Values()
	if len(retrievedValues) != 3 {
		t.Errorf("expected 3 values, got %d", len(retrievedValues))
	}
}

func TestValuesEmpty(t *testing.T) {
	header := NewHeader("X-Custom", []*Value{})

	values := header.Values()
	if len(values) != 0 {
		t.Errorf("expected 0 values, got %d", len(values))
	}
}

func TestAddValue(t *testing.T) {
	header := NewHeader("Accept", []*Value{})
	value := &Value{}

	header.AddValue(value)

	if len(header.Values()) != 1 {
		t.Errorf("expected 1 value after AddValue, got %d", len(header.Values()))
	}
}

func TestAddValueMultiple(t *testing.T) {
	header := NewHeader("Accept", []*Value{})

	header.AddValue(&Value{})
	header.AddValue(&Value{})
	header.AddValue(&Value{})

	if len(header.Values()) != 3 {
		t.Errorf("expected 3 values after multiple AddValue calls, got %d", len(header.Values()))
	}
}

func TestAddValueToExistingValues(t *testing.T) {
	initialValues := []*Value{{}, {}}
	header := NewHeader("Accept", initialValues)

	header.AddValue(&Value{})

	if len(header.Values()) != 3 {
		t.Errorf("expected 3 values, got %d", len(header.Values()))
	}
}

func TestAddValues(t *testing.T) {
	header := NewHeader("Accept", []*Value{})
	newValues := []Value{{}, {}, {}}

	header.AddValues(newValues)

	if len(header.Values()) != 3 {
		t.Errorf("expected 3 values after AddValues, got %d", len(header.Values()))
	}
}

func TestAddValuesMultipleCalls(t *testing.T) {
	header := NewHeader("Accept", []*Value{})

	header.AddValues([]Value{{}, {}})
	header.AddValues([]Value{{}, {}})

	if len(header.Values()) != 4 {
		t.Errorf("expected 4 values after multiple AddValues calls, got %d", len(header.Values()))
	}
}

func TestAddValuesToExistingValues(t *testing.T) {
	initialValues := []*Value{{}}
	header := NewHeader("Accept", initialValues)
	newValues := []Value{{}, {}}

	header.AddValues(newValues)

	if len(header.Values()) != 3 {
		t.Errorf("expected 3 values total, got %d", len(header.Values()))
	}
}

func TestAddValuesEmpty(t *testing.T) {
	header := NewHeader("Accept", []*Value{})

	header.AddValues([]Value{})

	if len(header.Values()) != 0 {
		t.Errorf("expected 0 values after adding empty slice, got %d", len(header.Values()))
	}
}

func TestHeaderFromStringComplexContentType(t *testing.T) {
	header, err := NewHeaderFromString("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")

	if err != nil {
		t.Errorf("NewHeaderFromString() error = %v", err)
	}
	if header == nil {
		t.Fatal("NewHeaderFromString() returned nil header, expected non-nil header")
	}
	if len(header.Values()) != 1 {
		t.Errorf("expected exactly 1 value, got %d", len(header.Values()))
	}
}

func TestHeaderFromStringComplexAccept(t *testing.T) {
	header, err := NewHeaderFromString("Accept", "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8")

	if err != nil {
		t.Errorf("NewHeaderFromString() error = %v", err)
	}
	if header == nil {
		t.Fatal("NewHeaderFromString() returned nil header, expected non-nil header")
	}
	if len(header.Values()) != 4 {
		t.Errorf("expected exactly 4 values, got %d", len(header.Values()))
	}
}

func TestNewHeaderPreservesName(t *testing.T) {
	testCases := []string{
		"Content-Type",
		"Authorization",
		"X-Custom-Header",
		"Accept-Encoding",
	}

	for _, name := range testCases {
		header := NewHeader(name, []*Value{})
		if header.Name() != name {
			t.Errorf("expected name '%s', got '%s'", name, header.Name())
		}
	}
}

func TestNewHeaderFromStringPreservesName(t *testing.T) {
	testCases := []string{
		"Content-Type",
		"Authorization",
		"X-Custom-Header",
	}

	for _, name := range testCases {
		header, err := NewHeaderFromString(name, "value")
		if err != nil {
			t.Errorf("NewHeaderFromString() error = %v", err)
		}
		if header.Name() != name {
			t.Errorf("expected name '%s', got '%s'", name, header.Name())
		}
	}
}
