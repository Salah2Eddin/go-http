package response

import (
	"bytes"
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
)

func TestNewEmptyResponse(t *testing.T) {
	line := &StatusLine{Version: "HTTP/1.1", Code: 200, Phrase: "OK"}
	resp := NewEmptyResponse(line)

	if resp.Line != line {
		t.Errorf("Line = %v, want %v", resp.Line, line)
	}
	if resp.Headers == nil {
		t.Error("NewEmptyResponse() returned response with nil Headers, expected non-nil initialized headers")
	}
	if len(resp.Body) != 0 {
		t.Errorf("Body = %v, want empty", resp.Body)
	}
}

func TestNewResponse(t *testing.T) {
	line := &StatusLine{Version: "HTTP/1.1", Code: 404, Phrase: "Not Found"}
	headers := httpheaders.New()
	body := []byte("hello")

	resp := NewResponse(line, headers, body)

	if resp.Line != line {
		t.Errorf("Line = %v, want %v", resp.Line, line)
	}
	if resp.Headers != headers {
		t.Errorf("Headers = %v, want %v", resp.Headers, headers)
	}
	if !bytes.Equal(resp.Body, body) {
		t.Errorf("Body = %v, want %v", resp.Body, body)
	}
}
