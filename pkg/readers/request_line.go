package readers

import (
	"bufio"
)

type requestLineReader struct {
}

func (r requestLineReader) Read(reader *bufio.Reader) ([]byte, error) {
	var buf []byte
	err := readLine(reader, &buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
