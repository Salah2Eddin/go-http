package httpbody

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/httpheader"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"strconv"
)

const (
	contentLengthHeaderName = "content-length"
)

func GetRequestBody(reader *bufio.Reader, headers httpheader.Headers) (*[]byte, error) {
	lengthHeader, exists := headers.Get(contentLengthHeaderName)
	if !exists {
		return &[]byte{}, nil
	}

	lengthString := lengthHeader.Values()[0].Value()
	length, err := strconv.Atoi(lengthString)
	if err != nil {
		return nil, &pkgerrors.ErrInvalidContentLength{Length: lengthString}
	}

	// TODO: Handle chunked content here
	var bodyReader IBodyReader = nil
	bodyReader = BodyReader{}

	body := make([]byte, length)
	err = bodyReader.read(reader, &body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
