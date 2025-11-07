package serializers

import (
	"bytes"
	"strconv"

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
	// Serialize body first to know its final length after any transfer encoding
	bodySerializer := bodySerializerFactory(resp)
	tmpBodyBuf := &bytes.Buffer{}
	bodySerializer.Serialize(resp.Body, tmpBodyBuf)

	// Decide on content-length: only when no transfer-encoding and header not preset
	_, hasTransferEncoding := resp.Headers.Get("transfer-encoding")
	_, hasLength := resp.Headers.Get("content-length")
	if !hasTransferEncoding && !hasLength {
		resp.Headers.AddFromString("content-length", strconv.Itoa(tmpBodyBuf.Len()))
	}

	r.statusSerializer.Serialize(resp.Line, buf)
	r.headersSerializer.Serialize(resp.Headers, buf)

	// Empty line between headers and body
	buf.WriteString("\r\n")

	// Write finalized body bytes
	buf.Write(tmpBodyBuf.Bytes())
}

type StatusLineSerializer struct{}

func (StatusLineSerializer) Serialize(status *response.StatusLine, buf *bytes.Buffer) {
	buf.WriteString(status.Version)
	buf.WriteByte(' ')
	buf.WriteString(strconv.Itoa(status.Code))
	buf.WriteByte(' ')
	buf.WriteString(status.Phrase)
	buf.WriteString("\r\n")
}
