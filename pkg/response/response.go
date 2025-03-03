package response

type Response struct {
	Line    *StatusLine
	Headers *ResponseHeaders
	Body    *[]byte
}

func NewEmptyResponse(line *StatusLine) *Response {
	body := make([]byte, 0)
	return &Response{
		Line:    line,
		Headers: &ResponseHeaders{},
		Body:    &body,
	}
}

func NewResponse(line *StatusLine, headers *ResponseHeaders, body *[]byte) *Response {
	return &Response{
		Line:    line,
		Headers: headers,
		Body:    body,
	}
}

func (res Response) String() string {
	line_bytes := res.Line.String()
	header_bytes := res.Headers.String()

	return line_bytes + header_bytes + "\r\n" + string(*res.Body)
}

func (res Response) Bytes() []byte {
	line_bytes := []byte(res.Line.String())
	header_bytes := []byte(res.Headers.String())

	bytes := make([]byte, 0)
	bytes = append(bytes, line_bytes...)
	bytes = append(bytes, header_bytes...)

	// empty line between headers and body
	bytes = append(bytes, []byte("\r\n")...)

	bytes = append(bytes, *res.Body...)

	return bytes
}
