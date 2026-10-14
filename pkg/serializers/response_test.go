package serializers

import (
	"bytes"
	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"strings"
	"testing"
)

func TestResponseSerializerSerialize(t *testing.T) {
	tests := []struct {
		name                string
		body                []byte
		hasTransferEncoding bool
		hasContentLength    bool
		expectedHasLength   bool
	}{
		{"no headers", []byte("test"), false, false, true},
		{"with transfer encoding", []byte("test"), true, false, false},
		{"content length already set", []byte("test"), false, true, false},
		{"empty body", []byte{}, false, false, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			headers := httpheaders.New()
			if tc.hasTransferEncoding {
				err := headers.AddFromString("transfer-encoding", "chunked")
				if err != nil {
					t.Fatalf("AddFromString failed: %v", err)
				}
			}
			if tc.hasContentLength {
				err := headers.AddFromString("content-length", "5")
				if err != nil {
					t.Fatalf("AddFromString failed: %v", err)
				}
			}
			statusLine := &response.StatusLine{
				Version: "HTTP/1.1",
				Code:    200,
				Phrase:  "OK",
			}
			resp := &response.Response{
				Line:    statusLine,
				Headers: headers,
				Body:    tc.body,
			}
			serializer := NewResponseSerializer()
			buf := &bytes.Buffer{}
			serializer.Serialize(resp, buf)
			result := buf.String()

			if tc.expectedHasLength {
				if !strings.Contains(result, "content-length:") {
					t.Errorf("Serialize() missing content-length header")
				}
			} else {
				contentLengthCount := strings.Count(result, "content-length:")
				if tc.hasContentLength && contentLengthCount != 1 {
					t.Errorf("Serialize() added duplicate content-length")
				}
			}

			if !strings.Contains(result, "\r\n\r\n") {
				t.Errorf("Serialize() missing blank line between headers and body")
			}

			if !strings.HasSuffix(result, string(tc.body)) {
				t.Errorf("Serialize() body not at end")
			}
		})
	}
}

func TestStatusLineSerializerSerialize(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		code     int
		phrase   string
		expected string
	}{
		{"ok", "HTTP/1.1", 200, "OK", "HTTP/1.1 200 OK\r\n"},
		{"not found", "HTTP/1.1", 404, "Not Found", "HTTP/1.1 404 Not Found\r\n"},
		{"server error", "HTTP/1.0", 500, "Internal Server Error", "HTTP/1.0 500 Internal Server Error\r\n"},
		{"continue", "HTTP/1.1", 100, "Continue", "HTTP/1.1 100 Continue\r\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status := &response.StatusLine{
				Version: tc.version,
				Code:    tc.code,
				Phrase:  tc.phrase,
			}
			serializer := StatusLineSerializer{}
			buf := &bytes.Buffer{}
			serializer.Serialize(status, buf)
			result := buf.String()
			if result != tc.expected {
				t.Errorf("Serialize() = %q, want %q", result, tc.expected)
			}
		})
	}
}
