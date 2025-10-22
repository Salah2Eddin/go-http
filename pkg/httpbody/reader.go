package httpbody

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"io"
)

type IBodyReader interface {
	read(reader *bufio.Reader, buf *[]byte) error
}

type BodyReader struct {
}

func (r BodyReader) read(reader *bufio.Reader, buf *[]byte) error {
	_, err := io.ReadFull(reader, *buf)
	if err != nil {
		return &pkgerrors.ErrIncorrectContentLength{}
	}
	return nil
}
