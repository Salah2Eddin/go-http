package request

import (
	"github.com/Salah2Eddin/go-http/pkg/httpheader"
	"github.com/Salah2Eddin/go-http/pkg/reqline"
	"github.com/Salah2Eddin/go-http/pkg/uri"
)

type Request struct {
	line    *reqline.RequestLine
	headers *httpheader.Headers
	Body    *[]byte
}

func NewRequest(line *reqline.RequestLine, headers *httpheader.Headers, body *[]byte) Request {
	return Request{
		line:    line,
		headers: headers,
		Body:    body,
	}
}

func (req *Request) Uri() *uri.Uri {
	return req.line.Uri
}

func (req *Request) GetUriParameter(param string) (string, bool) {
	return req.line.Uri.GetQueryParameter(param)
}

func (req *Request) Method() string {
	return req.line.Method
}

func (req *Request) Version() string {
	return req.line.Version
}

func (req *Request) GetHeader(name string) (*httpheader.Header, bool) {
	// make it case in-sensitive
	return req.headers.Get(name)
}
