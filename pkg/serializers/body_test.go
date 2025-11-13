package serializers

import (
	"bytes"
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
	"github.com/Salah2Eddin/go-http/pkg/response"
)

func TestBodySerializerFactory(t *testing.T) {
	tests := []struct {
		name         string
		encoding     string
		hasHeader    bool
		expectedType string
	}{
		{"no transfer encoding", "", false, "default"},
		{"unknown encoding", "gzip", true, "default"},
		{"default case", "chunked", true, "default"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			headers := httpheaders.New()
			if tc.hasHeader {
				err := headers.AddFromString("transfer-encoding", tc.encoding)
				if err != nil {
					t.Fatalf("AddFromString failed: %v", err)
				}
			}
			resp := &response.Response{Headers: headers}
			result := bodySerializerFactory(resp)
			if _, ok := result.(BodySerializer); !ok {
				t.Errorf("bodySerializerFactory() returned %T, want BodySerializer", result)
			}
		})
	}
}

func TestBodySerializerSerialize(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected []byte
	}{
		{"empty body", []byte{}, []byte{}},
		{"simple body", []byte("hello"), []byte("hello")},
		{"binary data", []byte{0x00, 0x01, 0x02}, []byte{0x00, 0x01, 0x02}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			serializer := BodySerializer{}
			buf := &bytes.Buffer{}
			serializer.Serialize(tc.body, buf)
			result := buf.Bytes()
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("Serialize() = %v, want %v", result, tc.expected)
			}
		})
	}
}
