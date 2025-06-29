package request

import (
	"github.com/Salah2Eddin/go-http/pkg/uri"
)

type Request struct {
	line    Line
	headers Headers
	Body    *[]byte
}

func NewRequest(line Line, headers Headers, body *[]byte) Request {
	return Request{
		line:    line,
		headers: headers,
		Body:    body,
	}
}

func (req *Request) Uri() uri.Uri {
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

func (req *Request) GetHeader(name string) (Header, bool) {
	// make it case in-sensitive
	return req.headers.Get(name)
}
