package response

type Response struct {
	Line    *StatusLine
	Headers *ResponseHeaders
	Body    *[]byte
}

func NewResponse(line *StatusLine, headers *ResponseHeaders, body *[]byte) *Response {
	return &Response{
		Line:    line,
		Headers: headers,
		Body:    body,
	}
}

func (res Response) Bytes() []byte {
	line_bytes := []byte(res.Line.String())
	header_bytes := []byte(res.Headers.String())

	bytes := make([]byte, len(line_bytes)+len(header_bytes)+len(*res.Body))
	bytes = append(bytes, line_bytes...)
	bytes = append(bytes, header_bytes...)
	bytes = append(bytes, *res.Body...)

	return bytes
}
