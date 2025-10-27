package response

import (
	"github.com/Salah2Eddin/go-http/pkg/httpheader"
)

type Response struct {
	Line    *StatusLine
	Headers *httpheader.Headers
	Body    *[]byte
}

func NewEmptyResponse(line *StatusLine) *Response {
	body := make([]byte, 0)
	return &Response{
		Line:    line,
		Headers: &httpheader.Headers{},
		Body:    &body,
	}
}

func NewResponse(line *StatusLine, headers *httpheader.Headers, body *[]byte) *Response {
	return &Response{
		Line:    line,
		Headers: headers,
		Body:    body,
	}
}
