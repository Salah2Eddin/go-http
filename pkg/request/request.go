package request

type Request struct {
	line    *RequestLine
	headers *RequestHeaders
	Body    *[]byte
}

func NewRequest(line *RequestLine, headers *RequestHeaders, body *[]byte) *Request {
	return &Request{
		line:    line,
		headers: headers,
		Body:    body,
	}
}

func (req *Request) Uri() string {
	return req.line.Uri.String()
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
