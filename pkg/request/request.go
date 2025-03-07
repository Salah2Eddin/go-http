package request

import "ducky/http/pkg/uri"

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
	return req.line.Uri.GetParameter(param)
}

func (req *Request) Method() string {
	return req.line.Method
}

func (req *Request) Version() string {
	return req.line.Version
}

func (req *Request) GetHeader(header string) ([]string, bool) {
	return req.headers.Get(header)
}
