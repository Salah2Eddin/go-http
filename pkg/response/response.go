package response

import "github.com/Salah2Eddin/go-http/pkg/httpheader"

type Response struct {
	Line    StatusLine
	Headers httpheader.Headers
	Body    *[]byte
}

func NewEmptyResponse(line StatusLine) Response {
	body := make([]byte, 0)
	return Response{
		Line:    line,
		Headers: httpheader.Headers{},
		Body:    &body,
	}
}

func NewResponse(line StatusLine, headers httpheader.Headers, body *[]byte) Response {
	return Response{
		Line:    line,
		Headers: headers,
		Body:    body,
	}
}

func (res *Response) String() string {
	lineBytes := res.Line.String()
	headerBytes := res.Headers.Serialize()

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
