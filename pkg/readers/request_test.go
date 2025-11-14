package readers

import (
	"bufio"
	"strings"
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
)

// bodyReaderFactory tests
func TestBodyReaderFactoryWithContentLength(t *testing.T) {
	r := NewRequestReader()
	headers := httpheaders.New()
	headers.AddFromString("Content-Length", "100")

	bodyReader, err := r.bodyReaderFactory(headers)

	if err != nil {
		t.Errorf("bodyReaderFactory() error = %v, expected nil", err)
	}
	if bodyReader == nil {
		t.Error("bodyReaderFactory() returned nil reader, expected non-nil EmptyBodyReader")
	}
	_, ok := bodyReader.(lengthBodyReader)
	if !ok {
		t.Errorf("bodyReaderFactory() returned %T, expected lengthBodyReader", bodyReader)
	}
}

func TestBodyReaderFactoryWithTransferEncoding(t *testing.T) {
	r := NewRequestReader()
	headers := httpheaders.New()
	headers.AddFromString("Transfer-Encoding", "chunked")

	bodyReader, err := r.bodyReaderFactory(headers)

	if err == nil {
		t.Error("bodyReaderFactory() returned nil error for unsupported transfer encoding, expected non-nil error")
	}
	_, ok := err.(pkgerrors.ErrUnsupportedBodyTransferEncoding)
	if !ok {
		t.Errorf("bodyReaderFactory() error type = %T, expected ErrUnsupportedBodyTransferEncoding", err)
	}
	if bodyReader != nil {
		t.Errorf("bodyReaderFactory() returned %v, expected nil", bodyReader)
	}
}

func TestBodyReaderFactoryWithoutContentLengthOrTransferEncoding(t *testing.T) {
	r := NewRequestReader()
	headers := httpheaders.New()

	bodyReader, err := r.bodyReaderFactory(headers)

	if err != nil {
		t.Errorf("bodyReaderFactory() error = %v, expected nil", err)
	}
	if bodyReader == nil {
		t.Error("bodyReaderFactory() returned nil reader, expected non-nil EmptyBodyReader")
	}
	_, ok := bodyReader.(EmptyBodyReader)
	if !ok {
		t.Errorf("bodyReaderFactory() returned %T, expected EmptyBodyReader", bodyReader)
	}
}

func TestBodyReaderFactoryInvalidContentLength(t *testing.T) {
	r := NewRequestReader()
	headers := httpheaders.New()
	headers.AddFromString("Content-Length", "not-a-number")

	bodyReader, err := r.bodyReaderFactory(headers)

	if err == nil {
		t.Error("bodyReaderFactory() returned nil error for unsupported transfer encoding, expected non-nil error")
	}
	_, ok := err.(pkgerrors.ErrInvalidContentLength)
	if !ok {
		t.Errorf("bodyReaderFactory() error type = %T, expected ErrInvalidContentLength", err)
	}
	if bodyReader != nil {
		t.Errorf("bodyReaderFactory() returned %v, expected nil", bodyReader)
	}
}

func TestBodyReaderFactoryTransferEncodingTakesPrecedence(t *testing.T) {
	r := NewRequestReader()
	headers := httpheaders.New()
	headers.AddFromString("Transfer-Encoding", "chunked")
	headers.AddFromString("Content-Length", "100")

	bodyReader, err := r.bodyReaderFactory(headers)

	if err == nil {
		t.Error("bodyReaderFactory() returned nil error for unsupported transfer encoding, expected non-nil error")
	}
	_, ok := err.(pkgerrors.ErrUnsupportedBodyTransferEncoding)
	if !ok {
		t.Errorf("bodyReaderFactory() error type = %T, expected ErrUnsupportedBodyTransferEncoding", err)
	}
	if bodyReader != nil {
		t.Errorf("bodyReaderFactory() returned %v, expected nil", bodyReader)
	}
}

// Parse tests
func TestParseValidRequest(t *testing.T) {
	data := "GET /path HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 0\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	r := NewRequestReader()

	result, err := r.Parse(reader)

	if err != nil {
		t.Errorf("Parse() error = %v, expected nil", err)
	}
	if result == nil {
		t.Error("Parse() returned nil request")
	}
}

func TestParseInvalidRequestLine(t *testing.T) {
	data := ""
	reader := bufio.NewReader(strings.NewReader(data))
	r := NewRequestReader()

	result, err := r.Parse(reader)

	if err == nil {
		t.Error("Parse() returned nil error, expected non-nil error")
	}
	if result != nil {
		t.Errorf("Parse() returned %v, expected nil", result)
	}
}

func TestParseInvalidHeaders(t *testing.T) {
	data := "GET /path HTTP/1.1\r\nInvalid\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	r := NewRequestReader()

	result, err := r.Parse(reader)

	if err == nil {
		t.Error("Parse() returned nil error, expected non-nil error")
	}
	if result != nil {
		t.Errorf("Parse() returned %v, expected nil", result)
	}
}

func TestParseUnsupportedTransferEncoding(t *testing.T) {
	data := "GET /path HTTP/1.1\r\nTransfer-Encoding: chunked\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(data))
	r := NewRequestReader()

	result, err := r.Parse(reader)

	if err == nil {
		t.Error("Parse() returned nil error, expected non-nil error")
	}
	if result != nil {
		t.Errorf("Parse() returned %v, expected nil", result)
	}
}
