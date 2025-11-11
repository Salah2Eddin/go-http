package readers

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
)

// EmptyBodyReader tests
func TestEmptyBodyReaderWithEmptyBuffer(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte{}))
	e := EmptyBodyReader{}

	result, err := e.Read(reader)

	if err != nil {
		t.Errorf("Read() error = %v, expected nil", err)
	}
	if len(result) != 0 {
		t.Errorf("Read() returned %d bytes, expected 0", len(result))
	}
}

func TestEmptyBodyReaderWithBufferedData(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("unexpected data"))
	e := EmptyBodyReader{}

	result, err := e.Read(reader)

	if err == nil {
		t.Error("Read() expected error, got nil")
	}
	var errExpectedEmptyBody pkgerrors.ErrExpectedEmptyBody
	if !errors.As(err, &errExpectedEmptyBody) {
		t.Errorf("Read() error type = %T, expected ErrExpectedEmptyBody", err)
	}
	if result != nil {
		t.Errorf("Read() returned %v, expected nil", result)
	}
}

func TestEmptyBodyReaderWithSingleByte(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("a"))
	e := EmptyBodyReader{}

	result, err := e.Read(reader)

	if err == nil {
		t.Error("Read() expected error, got nil")
	}
	var errExpectedEmptyBody pkgerrors.ErrExpectedEmptyBody
	if !errors.As(err, &errExpectedEmptyBody) {
		t.Errorf("Read() error type = %T, expected ErrExpectedEmptyBody", err)
	}
	if result != nil {
		t.Errorf("Read() returned %v, expected nil", result)
	}
}

func TestEmptyBodyReaderReturnsEmptySlice(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte{}))
	e := EmptyBodyReader{}

	result, err := e.Read(reader)

	if err != nil {
		t.Errorf("Read() error = %v", err)
	}
	if result == nil {
		t.Error("Read() returned nil, expected empty slice")
	}
	if len(result) != 0 {
		t.Errorf("Read() length = %d, expected 0", len(result))
	}
}

// newLengthBodyReader tests
func TestNewLengthBodyReader(t *testing.T) {
	length := 10
	reader := newLengthBodyReader(length)

	if reader.length != length {
		t.Errorf("newLengthBodyReader() length = %d, expected %d", reader.length, length)
	}
}

func TestNewLengthBodyReaderZeroLength(t *testing.T) {
	reader := newLengthBodyReader(0)

	if reader.length != 0 {
		t.Errorf("newLengthBodyReader() length = %d, expected 0", reader.length)
	}
}

// lengthBodyReader.Read tests
func TestLengthBodyReaderReadExactLength(t *testing.T) {
	data := "Hello World"
	reader := bufio.NewReader(strings.NewReader(data))
	lr := newLengthBodyReader(len(data))

	result, err := lr.Read(reader)

	if err != nil {
		t.Errorf("Read() error = %v, expected nil", err)
	}
	if string(result) != data {
		t.Errorf("Read() returned %q, expected %q", string(result), data)
	}
	if len(result) != len(data) {
		t.Errorf("Read() length = %d, expected %d", len(result), len(data))
	}
}

func TestLengthBodyReaderReadInsufficientData(t *testing.T) {
	data := "Hello"
	reader := bufio.NewReader(strings.NewReader(data))
	lr := newLengthBodyReader(20)

	result, err := lr.Read(reader)

	if err == nil {
		t.Error("Read() expected error, got nil")
	}
	var errIncorrectContentLength pkgerrors.ErrIncorrectContentLength
	if !errors.As(err, &errIncorrectContentLength) {
		t.Errorf("Read() error type = %T, expected ErrIncorrectContentLength", err)
	}
	if result != nil {
		t.Errorf("Read() returned %v, expected nil", result)
	}
}
