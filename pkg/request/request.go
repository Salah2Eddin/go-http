package request

type Request struct {
	Line    *RequestLine
	Headers *RequestHeaders
	Body    *[]byte
}

func NewRequest(line *RequestLine, headers *RequestHeaders, body *[]byte) *Request {
	return &Request{
		Line:    line,
		Headers: headers,
		Body:    body,
	}
}
