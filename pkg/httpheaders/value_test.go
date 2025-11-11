package httpheaders

import (
	"testing"
)

func TestNewHeaderValues(t *testing.T) {
	values := NewHeaderValues("application/json")
	if values == nil {
		t.Fatal("NewHeaderValues() returned nil")
	}
	if len(values) == 0 {
		t.Error("expected at least one value")
	}
}

func TestNewHeaderValuesMultiple(t *testing.T) {
	values := NewHeaderValues("application/json, text/html, text/plain")
	if values == nil {
		t.Fatal("NewHeaderValues() returned nil")
	}
	if len(values) < 3 {
		t.Errorf("expected at least 3 values, got %d", len(values))
	}
}

func TestNewHeaderValuesWithParameters(t *testing.T) {
	values := NewHeaderValues("text/html; charset=utf-8")
	if values == nil {
		t.Fatal("NewHeaderValues() returned nil")
	}
	if len(values) == 0 {
		t.Error("expected at least one value")
	}
}

func TestNewHeaderValuesEmpty(t *testing.T) {
	values := NewHeaderValues("")
	if values != nil {
		t.Fatal("expected NewHeaderValues() to return nil for empty string")
	}
}

func TestNewHeaderValueFromBytes(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("application/json"), []byte(""))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}
	if value.Value() != "application/json" {
		t.Errorf("expected 'application/json', got '%s'", value.Value())
	}
}

func TestNewHeaderValueFromBytesWithOneParam(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/html"), []byte("charset=utf-8"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}
	if value.Value() != "text/html" {
		t.Errorf("expected 'text/html', got '%s'", value.Value())
	}

	charSet, exists := value.GetParam("charset")
	if !exists {
		t.Error("expected 'charset' parameter to exist")
	}
	if charSet != "utf-8" {
		t.Errorf("expected 'utf-8', got '%s'", charSet)
	}
}

func TestNewHeaderValueFromBytesWithMultipleParams(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("multipart/form-data"), []byte("boundary=----WebKit;charset=utf-8"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	boundary, exists := value.GetParam("boundary")
	if !exists {
		t.Error("expected 'boundary' parameter to exist")
	}
	if boundary != "----WebKit" {
		t.Errorf("expected '----WebKit', got '%s'", boundary)
	}

	charSet, exists := value.GetParam("charset")
	if !exists {
		t.Error("expected 'charset' parameter to exist")
	}
	if charSet != "utf-8" {
		t.Errorf("expected 'utf-8', got '%s'", charSet)
	}
}

func TestNewHeaderValueFromBytesNoParams(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("application/json"), []byte(""))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	params := value.Params()
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestNewHeaderValueFromBytesInvalidParams(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/plain"), []byte("invalid;malformed"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	params := value.Params()
	if len(params) != 0 {
		t.Errorf("expected 0 params for invalid format, got %d", len(params))
	}
}

func TestValue(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("application/json"), []byte(""))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	if value.Value() != "application/json" {
		t.Errorf("expected 'application/json', got '%s'", value.Value())
	}
}

func TestParams(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/html"), []byte("charset=utf-8;format=flowed"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	params := value.Params()
	if len(params) < 2 {
		t.Errorf("expected at least 2 params, got %d", len(params))
	}
}

func TestParamsEmpty(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("application/json"), []byte(""))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	params := value.Params()
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestSetParam(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/plain"), []byte(""))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	value.SetParam("charset", "iso-8859-1")

	retrieved, exists := value.GetParam("charset")
	if !exists {
		t.Error("expected 'charset' parameter to exist after SetParam")
	}
	if retrieved != "iso-8859-1" {
		t.Errorf("expected 'iso-8859-1', got '%s'", retrieved)
	}
}

func TestSetParamOverwrite(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/html"), []byte("charset=utf-8"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	value.SetParam("charset", "iso-8859-1")

	retrieved, exists := value.GetParam("charset")
	if !exists {
		t.Error("expected 'charset' parameter to exist")
	}
	if retrieved != "iso-8859-1" {
		t.Errorf("expected 'iso-8859-1', got '%s'", retrieved)
	}
}

func TestSetParamMultiple(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("multipart/form-data"), []byte(""))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	value.SetParam("boundary", "----WebKit")
	value.SetParam("charset", "utf-8")

	params := value.Params()
	if len(params) != 2 {
		t.Errorf("expected 2 params, got %d", len(params))
	}
}

func TestGetParam(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/html"), []byte("charset=utf-8"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	retrieved, exists := value.GetParam("charset")
	if !exists {
		t.Error("expected 'charset' parameter to exist")
	}
	if retrieved != "utf-8" {
		t.Errorf("expected 'utf-8', got '%s'", retrieved)
	}
}

func TestGetParamNotFound(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/html"), []byte("charset=utf-8"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	_, exists := value.GetParam("format")
	if exists {
		t.Error("expected 'format' parameter to not exist")
	}
}

func TestGetParamEmpty(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("application/json"), []byte(""))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	_, exists := value.GetParam("charset")
	if exists {
		t.Error("expected parameter to not exist in empty params")
	}
}

func TestDeleteParam(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/html"), []byte("charset=utf-8"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	value.DeleteParam("charset")

	_, exists := value.GetParam("charset")
	if exists {
		t.Error("expected 'charset' parameter to be deleted")
	}
}

func TestDeleteParamMultiple(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/html"), []byte("charset=utf-8;format=flowed"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	value.DeleteParam("charset")

	_, exists := value.GetParam("charset")
	if exists {
		t.Error("expected 'charset' parameter to be deleted")
	}

	format, exists := value.GetParam("format")
	if !exists {
		t.Error("expected 'format' parameter to still exist")
	}
	if format != "flowed" {
		t.Errorf("expected 'flowed', got '%s'", format)
	}
}

func TestDeleteParamNotExists(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("application/json"), []byte(""))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	value.DeleteParam("charset")

	params := value.Params()
	if len(params) != 0 {
		t.Errorf("expected 0 params after delete, got %d", len(params))
	}
}

func TestComplexParamOperations(t *testing.T) {
	value := NewHeaderValueFromBytes([]byte("text/html"), []byte("charset=utf-8;format=flowed"))
	if value == nil {
		t.Fatal("NewHeaderValueFromBytes() returned nil")
	}

	value.SetParam("level", "1")
	value.DeleteParam("format")
	value.SetParam("charset", "iso-8859-1")

	charset, _ := value.GetParam("charset")
	if charset != "iso-8859-1" {
		t.Errorf("expected 'iso-8859-1', got '%s'", charset)
	}

	_, formatExists := value.GetParam("format")
	if formatExists {
		t.Error("expected 'format' to be deleted")
	}

	level, _ := value.GetParam("level")
	if level != "1" {
		t.Errorf("expected 'level' to be '1', got '%s'", level)
	}
}
