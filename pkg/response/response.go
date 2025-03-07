package response

type Response struct {
	Line    StatusLine
	Headers Headers
	Body    *[]byte
}

func NewEmptyResponse(line StatusLine) Response {
	body := make([]byte, 0)
	return Response{
		Line:    line,
		Headers: Headers{},
		Body:    &body,
	}
}

func NewResponse(line StatusLine, headers Headers, body *[]byte) Response {
	return Response{
		Line:    line,
		Headers: headers,
		Body:    body,
	}
}

func (res *Response) String() string {
	lineBytes := res.Line.String()
	headerBytes := res.Headers.String()

	return lineBytes + headerBytes + "\r\n" + string(*res.Body)
}

func (res *Response) Bytes() []byte {
	lineBytes := res.Line.Bytes()
	headerBytes := res.Headers.Bytes()

	bytes := make([]byte, 0)
	bytes = append(bytes, lineBytes...)
	bytes = append(bytes, headerBytes...)

	// empty line between headers and body
	bytes = append(bytes, []byte("\r\n")...)

	bytes = append(bytes, *res.Body...)

	return bytes
}
