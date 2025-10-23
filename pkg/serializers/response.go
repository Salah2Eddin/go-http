package serializers

import (
	"bytes"
	"github.com/Salah2Eddin/go-http/pkg/response"
)

type ResponseSerializer struct {
	statusSerializer  StatusLineSerializer
	headersSerializer HeadersSerializer
}

func bodySerializerFactory(resp *response.Response) ISerializer[[]byte] {
	if _, ok := resp.Headers.Get("transfer-encoding"); ok {
		// TODO: handle cases like chunked body
		return BodySerializer{}
	}
	return BodySerializer{}
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
