package readers

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

// headersReader.Read tests
func TestHeadersReaderReadSingleHeader(t *testing.T) {
	data := "Content-Type: application/json\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	h := headersReader{}

	result, err := h.Read(reader)

	if err != nil {
		t.Errorf("Read() error = %v, expected nil", err)
	}
	if len(result) != 1 {
		t.Errorf("Read() count = %d, expected 1", len(result))
	}
	if string(result[0]) != "Content-Type: application/json" {
		t.Errorf("Read() first header = %q, expected %q", string(result[0]), "Content-Type: application/json")
	}
}

func TestHeadersReaderReadMultipleHeaders(t *testing.T) {
	data := "Content-Type: application/json\r\nAccept: text/html\r\nAuthorization: Bearer token\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	h := headersReader{}

	result, err := h.Read(reader)

	if err != nil {
		t.Errorf("Read() error = %v, expected nil", err)
	}
	if len(result) != 3 {
		t.Errorf("Read() count = %d, expected 3", len(result))
	}
	if string(result[0]) != "Content-Type: application/json" {
		t.Errorf("Read() first header = %q", string(result[0]))
	}
	if string(result[1]) != "Accept: text/html" {
		t.Errorf("Read() second header = %q", string(result[1]))
	}
	if string(result[2]) != "Authorization: Bearer token" {
		t.Errorf("Read() third header = %q", string(result[2]))
	}
}

func TestHeadersReaderReadStopsAtEmptyLine(t *testing.T) {
	data := "Header1: value1\r\n\r\nShouldNotBeRead"
	reader := bufio.NewReader(strings.NewReader(data))
	h := headersReader{}

	result, err := h.Read(reader)

	if err != nil {
		t.Errorf("Read() error = %v, expected nil", err)
	}
	if len(result) != 1 {
		t.Errorf("Read() count = %d, expected 1", len(result))
	}
	if string(result[0]) != "Header1: value1" {
		t.Errorf("Read() header = %q", string(result[0]))
	}
}

func TestHeadersReaderReadError(t *testing.T) {
	data := ""
	reader := bufio.NewReader(strings.NewReader(data))
	h := headersReader{}

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
