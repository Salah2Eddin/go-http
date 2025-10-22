package request

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/httpbody"
	"github.com/Salah2Eddin/go-http/pkg/httpheader"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/reqline"
	"strconv"
)

const (
	contentLengthHeaderName = "content-length"
)

func FromReader(reader *bufio.Reader) (Request, error) {
	var req Request

	err := req.readRequestLine(reader)
	if err != nil {
		return Request{}, err
	}

	err = req.readHeaders(reader)
	if err != nil {
		return Request{}, err
	}

	err = req.readBody(reader)
	if err != nil {
		return Request{}, err
	}
	return req, nil
}

func (req *Request) readRequestLine(reader *bufio.Reader) error {
	var buf []byte
	err := reqline.Read(reader, &buf)
	if err != nil {
		return err
	}

	req.line, err = reqline.ParseRequestLine(buf)
	return err
}

func (req *Request) readHeaders(reader *bufio.Reader) error {
	var buf *[][]byte
	err := httpheader.Read(reader, buf)
	if err != nil {
		return err
	}
	req.headers, err = httpheader.ParseRequestHeaders(buf)
	return err
}

func (req *Request) readBody(reader *bufio.Reader) error {
	lengthHeader, exists := req.headers.Get(contentLengthHeaderName)
	if !exists {
		return nil
	}

	lengthString := lengthHeader.Values()[0].Value()
	length, err := strconv.Atoi(lengthString)
	if err != nil {
		return &pkgerrors.ErrInvalidContentLength{Length: lengthString}
	}

	// TODO: Handle chunked content here
	var bodyReader httpbody.IBodyReader = nil
	bodyReader = httpbody.BodyReader{}

	buf := make([]byte, length)
	err = bodyReader.Read(reader, &buf)
	if err != nil {
		return err
	}

	req.Body = &buf
	return nil
}
