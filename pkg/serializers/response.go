package serializers

import (
	"bytes"
	"github.com/Salah2Eddin/go-http/pkg/response"
)

type ResponseSerializer struct {
	statusSerializer  StatusLineSerializer
	headersSerializer HeadersSerializer
}

func NewResponseSerializer() ResponseSerializer {
	return ResponseSerializer{
		statusSerializer:  StatusLineSerializer{},
		headersSerializer: HeadersSerializer{},
	}
}

func (r ResponseSerializer) Serialize(resp *response.Response, buf *bytes.Buffer) {
	r.statusSerializer.Serialize(&resp.Line, buf)
	r.headersSerializer.Serialize(&resp.Headers, buf)

	// Empty line between headers and body
	buf.WriteString("\r\n")

	bodySerializer := bodySerializerFactory(resp)
	bodySerializer.Serialize(resp.Body, buf)
}

type StatusLineSerializer struct{}

func (StatusLineSerializer) Serialize(status *response.StatusLine, buf *bytes.Buffer) {
	buf.WriteString(status.Version)
	buf.WriteByte(' ')
	buf.WriteString(status.Code)
	buf.WriteByte(' ')
	buf.WriteString(status.Phrase)
	buf.WriteString("\r\n")
}
