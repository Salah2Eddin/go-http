package readers

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
)

// checkCRLF tests
func TestCheckCRLF(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "valid CRLF",
			input:    []byte("line\r\n"),
			expected: true,
		},
		{
			name:     "only CR",
			input:    []byte("line\r"),
			expected: false,
		},
		{
			name:     "only LF",
			input:    []byte("line\n"),
			expected: false,
		},
		{
			name:     "reversed CRLF",
			input:    []byte("line\n\r"),
			expected: false,
		},
		{
			name:     "empty bytes",
			input:    []byte(""),
			expected: false,
		},
		{
			name:     "single byte",
			input:    []byte("a"),
			expected: false,
		},
		{
			name:     "exactly CRLF",
			input:    []byte("\r\n"),
			expected: true,
		},
		{
			name:     "multiple lines with CRLF",
			input:    []byte("line1\r\nline2\r\n"),
			expected: true,
		},
		{
			name:     "spaces before CRLF",
			input:    []byte("line   \r\n"),
			expected: true,
		},
		{
			name:     "no CRLF",
			input:    []byte("line"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := checkCRLF(tc.input)
			if result != tc.expected {
				t.Errorf("checkCRLF(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

// checkHeadersEnd tests
func TestCheckHeadersEnd(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "empty bytes",
			input:    []byte(""),
			expected: true,
		},
		{
			name:     "single byte",
			input:    []byte("a"),
			expected: false,
		},
		{
			name:     "single space",
			input:    []byte(" "),
			expected: false,
		},
		{
			name:     "CRLF",
			input:    []byte("\r\n"),
			expected: false,
		},
		{
			name:     "multiple bytes",
			input:    []byte("header"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := checkHeadersEnd(tc.input)
			if result != tc.expected {
				t.Errorf("checkHeadersEnd(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

// readLine tests
func TestReadLineSimple(t *testing.T) {
	data := "Hello World\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err != nil {
		t.Errorf("readLine() error = %v, expected nil", err)
	}
	if string(buf) != "Hello World" {
		t.Errorf("readLine() returned %q, expected %q", string(buf), "Hello World")
	}
	if len(buf) != 11 {
		t.Errorf("readLine() length = %d, expected 11", len(buf))
	}
}

func TestReadLineEmptyLine(t *testing.T) {
	data := "\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err != nil {
		t.Errorf("readLine() error = %v, expected nil", err)
	}
	if len(buf) != 0 {
		t.Errorf("readLine() length = %d, expected 0", len(buf))
	}
}

func TestReadLinePrefilledBuffer(t *testing.T) {
	data := "World\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := []byte("Hello ")

	err := readLine(reader, &buf)

	if err != nil {
		t.Errorf("readLine() error = %v, expected nil", err)
	}
	if string(buf) != "Hello World" {
		t.Errorf("readLine() returned %q, expected %q", string(buf), "Hello World")
	}
}

func TestReadLineRemovesCRLF(t *testing.T) {
	data := "test\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err != nil {
		t.Errorf("readLine() error = %v, expected nil", err)
	}
	// Verify CRLF was removed
	if checkCRLF(buf) {
		t.Error("readLine() did not remove CRLF")
	}
	if string(buf) != "test" {
		t.Errorf("readLine() returned %q, expected %q", string(buf), "test")
	}
}

func TestReadLineSpecialCharacters(t *testing.T) {
	data := "!@#$%^&*()\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err != nil {
		t.Errorf("readLine() error = %v, expected nil", err)
	}
	if string(buf) != "!@#$%^&*()" {
		t.Errorf("readLine() returned %q, expected %q", string(buf), "!@#$%^&*()")
	}
}

func TestReadLineWithSpaces(t *testing.T) {
	data := "  spaces  everywhere  \r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err != nil {
		t.Errorf("readLine() error = %v, expected nil", err)
	}
	if string(buf) != "  spaces  everywhere  " {
		t.Errorf("readLine() returned %q, expected %q", string(buf), "  spaces  everywhere  ")
	}
}

func TestReadLineLongLine(t *testing.T) {
	lineContent := strings.Repeat("a", 1000)
	data := lineContent + "\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err != nil {
		t.Errorf("readLine() error = %v, expected nil", err)
	}
	if string(buf) != lineContent {
		t.Errorf("readLine() length = %d, expected %d", len(buf), len(lineContent))
	}
}

func TestReadLineNoData(t *testing.T) {
	data := ""
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err == nil {
		t.Error("readLine() returned nil error, expected non-nil error (io.EOF)")
	}
	if err != io.EOF {
		t.Errorf("readLine() error = %v, expected io.EOF", err)
	}
}

func TestReadLineOnlyCR(t *testing.T) {
	data := "line\r"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err == nil {
		t.Error("readLine() returned nil error, expected non-nil error (io.EOF)")
	}
	if err != io.EOF {
		t.Errorf("readLine() error = %v, expected io.EOF", err)
	}
}

func TestReadLineOnlyLF(t *testing.T) {
	data := "line\n"
	reader := bufio.NewReader(strings.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err == nil {
		t.Error("readLine() returned nil error, expected non-nil error (io.EOF)")
	}
	if err != io.EOF {
		t.Errorf("readLine() error = %v, expected io.EOF", err)
	}
}

func TestReadLineMultipleLines(t *testing.T) {
	data := "line1\r\nline2\r\n"
	reader := bufio.NewReader(strings.NewReader(data))

	// Read first line
	buf1 := make([]byte, 0)
	err := readLine(reader, &buf1)
	if err != nil {
		t.Errorf("readLine() first line error = %v", err)
	}
	if string(buf1) != "line1" {
		t.Errorf("readLine() first line = %q, expected %q", string(buf1), "line1")
	}

	// Read second line
	buf2 := make([]byte, 0)
	err = readLine(reader, &buf2)
	if err != nil {
		t.Errorf("readLine() second line error = %v", err)
	}
	if string(buf2) != "line2" {
		t.Errorf("readLine() second line = %q, expected %q", string(buf2), "line2")
	}
}

func TestReadLineBinaryData(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, '\r', '\n'}
	reader := bufio.NewReader(bytes.NewReader(data))
	buf := make([]byte, 0)

	err := readLine(reader, &buf)

	if err != nil {
		t.Errorf("readLine() error = %v, expected nil", err)
	}
	if len(buf) != 5 {
		t.Errorf("readLine() length = %d, expected 5", len(buf))
	}
	if !bytes.Equal(buf, []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}) {
		t.Errorf("readLine() returned %v, expected %v", buf, []byte{0x00, 0x01, 0x02, 0xFF, 0xFE})
	}
}
