package readers

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"io"
)

type EmptyBodyReader struct {
}

func (e EmptyBodyReader) Read(reader *bufio.Reader) ([]byte, error) {
	return make([]byte, 0), nil
}

type lengthBodyReader struct {
	length int
}

func newLengthBodyReader(length int) lengthBodyReader {
	return lengthBodyReader{length: length}
}

func (r lengthBodyReader) Read(reader *bufio.Reader) ([]byte, error) {
	buf := make([]byte, r.length)
	n, err := io.ReadFull(reader, buf)
	if n != r.length || err != nil {
		return nil, &pkgerrors.ErrIncorrectContentLength{}
	}
	return buf, nil
}
