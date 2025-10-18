package httpbody

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/httpheader"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"io"
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

	body := make([]byte, length)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return nil, &pkgerrors.ErrIncorrectContentLength{Length: length}
	}

	return &body, nil
}
