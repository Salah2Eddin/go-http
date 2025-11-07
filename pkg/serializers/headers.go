package serializers

import (
	"bytes"
	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
)

type HeadersSerializer struct {
	headerSerializer HeaderSerializer
}

func (h HeadersSerializer) Serialize(headers *httpheaders.Headers, buf *bytes.Buffer) {
	for _, header := range headers.Headers() {
		h.headerSerializer.Serialize(header, buf)
		buf.WriteString("\r\n")
	}
	buf.WriteString("\r\n")
}

type HeaderSerializer struct {
}

func (HeaderSerializer) Serialize(header *httpheaders.Header, buf *bytes.Buffer) {
	buf.WriteString(header.Name())
	buf.WriteString(": ")

	values := header.Values()
	valueSerializer := ValueSerializer{}

	for i := range values {
		if i != 0 {
			buf.WriteString(headerValueSeparator)
		}
		formatter := formatterFactory(header.Name())
		valueSerializer.SetFormatter(formatter)
		valueSerializer.Serialize(values[i], buf)
	}
}

type ValueSerializer struct {
	formatter IValueFormatter
}

func (serializer *ValueSerializer) SetFormatter(formatter IValueFormatter) {
	serializer.formatter = formatter
}

func (serializer *ValueSerializer) Serialize(v *httpheaders.Value, buf *bytes.Buffer) {
	// Write main value
	buf.WriteString(serializer.formatter.format(v.Value()))

	// Write params
	for key, value := range v.Params() {
		buf.WriteString(headerParamSeparator)
		buf.WriteString(serializer.formatter.format(key))
		buf.WriteByte('=')
		buf.WriteString(serializer.formatter.format(value))
	}
}
