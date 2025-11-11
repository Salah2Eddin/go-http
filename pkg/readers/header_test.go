package readers

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

// headerReader.Read tests
func TestHeaderReaderReadSuccess(t *testing.T) {
	data := "Content-Type: application/json\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	h := headerReader{}

	result, err := h.Read(reader)

	if err != nil {
		t.Errorf("Read() error = %v, expected nil", err)
	}
	if string(result) != "Content-Type: application/json" {
		t.Errorf("Read() returned %q, expected %q", string(result), "Content-Type: application/json")
	}
}

func TestHeaderReaderReadError(t *testing.T) {
	data := ""
	reader := bufio.NewReader(strings.NewReader(data))
	h := headerReader{}

	result, err := h.Read(reader)

	if err == nil {
		t.Error("Read() expected error, got nil")
	}
	if err != io.EOF {
		t.Errorf("Read() error = %v, expected io.EOF", err)
	}
	if result != nil {
		t.Errorf("Read() returned %v, expected nil", result)
	}
}

func TestHeaderReaderReadMultiple(t *testing.T) {
	data := "Header1: value1\r\nHeader2: value2\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	h := headerReader{}

	result1, err := h.Read(reader)
	if err != nil {
		t.Errorf("Read() first error = %v", err)
	}
	if string(result1) != "Header1: value1" {
		t.Errorf("Read() first = %q, expected %q", string(result1), "Header1: value1")
	}

	result2, err := h.Read(reader)
	if err != nil {
		t.Errorf("Read() second error = %v", err)
	}
	if string(result2) != "Header2: value2" {
		t.Errorf("Read() second = %q, expected %q", string(result2), "Header2: value2")
	}
}
