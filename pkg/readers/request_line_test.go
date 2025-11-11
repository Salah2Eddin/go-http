package readers

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

// requestLineReader.Read tests
func TestRequestLineReaderReadSuccess(t *testing.T) {
	data := "GET /path HTTP/1.1\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	r := requestLineReader{}

	result, err := r.Read(reader)

	if err != nil {
		t.Errorf("Read() error = %v, expected nil", err)
	}
	if string(result) != "GET /path HTTP/1.1" {
		t.Errorf("Read() returned %q, expected %q", string(result), "GET /path HTTP/1.1")
	}
}

func TestRequestLineReaderReadError(t *testing.T) {
	data := ""
	reader := bufio.NewReader(strings.NewReader(data))
	r := requestLineReader{}

	result, err := r.Read(reader)

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
