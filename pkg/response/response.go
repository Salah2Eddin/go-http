package response

import (
	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
)

type Response struct {
	Line    *StatusLine
	Headers *httpheaders.Headers
	Body    *[]byte
}

func NewEmptyResponse(line *StatusLine) *Response {
	body := make([]byte, 0)
	return &Response{
		Line:    line,
		Headers: &httpheaders.Headers{},
		Body:    &body,
	}
}

func NewResponse(line *StatusLine, headers *httpheaders.Headers, body *[]byte) *Response {
	return &Response{
		Line:    line,
		Headers: headers,
		Body:    body,
	}
}
