package serializers

import (
	"bytes"
	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
	"testing"
)

func TestValueSerializerSetFormatter(t *testing.T) {
	serializer := &ValueSerializer{}
	formatter := IdentityFormatter{}
	serializer.SetFormatter(formatter)
	if serializer.formatter == nil {
		t.Errorf("SetFormatter failed to set formatter")
	}
}

func TestValueSerializerSerialize(t *testing.T) {
	tests := []struct {
		name      string
		valueStr  string
		formatter IValueFormatter
		expected  string
	}{
		{
			"simple value",
			"text/plain",
			IdentityFormatter{},
			"text/plain",
		},
		{
			"value with params",
			"text/plain;charset=utf-8",
			IdentityFormatter{},
			"text/plain;charset=utf-8",
		},
		{
			"value needing quotes",
			"hello world",
			QuotedFormatter{},
			`"hello world"`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			values := httpheaders.NewHeaderValues(tc.valueStr)
			if len(values) == 0 {
				t.Fatalf("NewHeaderValues(%q) returned no values", tc.valueStr)
			}
			serializer := &ValueSerializer{formatter: tc.formatter}
			buf := &bytes.Buffer{}
			serializer.Serialize(values[0], buf)
			result := buf.String()
			if result != tc.expected {
				t.Errorf("Serialize() = %q, expected %q", result, tc.expected)
			}
		})
	}
}

func TestHeaderSerializerSerialize(t *testing.T) {
	tests := []struct {
		name       string
		headerName string
		valueStr   string
		expected   string
	}{
		{
			"single value",
			"content-type",
			"text/plain",
			"content-type: text/plain",
		},
		{
			"multiple values",
			"accept",
			"text/plain,application/json",
			"accept: text/plain,application/json",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			values := httpheaders.NewHeaderValues(tc.valueStr)
			if len(values) == 0 {
				t.Fatalf("NewHeaderValues(%q) returned no values", tc.valueStr)
			}
			header := httpheaders.NewHeader(tc.headerName, values)
			serializer := HeaderSerializer{}
			buf := &bytes.Buffer{}
			serializer.Serialize(header, buf)
			result := buf.String()
			if result != tc.expected {
				t.Errorf("Serialize() = %q, expected %q", result, tc.expected)
			}
		})
	}
}

func TestHeadersSerializerSerialize(t *testing.T) {
	tests := []struct {
		name         string
		numHeaders   int
		expectedCRLF int
	}{
		{"single header", 1, 1},
		{"multiple headers", 3, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			headers := httpheaders.New()
			for i := 0; i < tc.numHeaders; i++ {
				err := headers.AddFromString("content-type", "text/plain")
				if err != nil {
					t.Fatalf("AddFromString failed: %v", err)
				}
			}
			serializer := HeadersSerializer{headerSerializer: HeaderSerializer{}}
			buf := &bytes.Buffer{}
			serializer.Serialize(headers, buf)
			valueCount := bytes.Count(buf.Bytes(), []byte("text/plain"))
			if valueCount != tc.expectedCRLF {
				t.Errorf("Serialize() wrote %d value, expected %d", valueCount, tc.expectedCRLF)
			}
		})
	}
}
