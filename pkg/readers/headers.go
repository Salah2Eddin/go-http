package readers

import (
	"bufio"
)

type headersReader struct {
	headerReader headerReader
}

func (h headersReader) Read(reader *bufio.Reader) ([][]byte, error) {
	buf := make([][]byte, 0)
	for {
		headerBuf, err := h.headerReader.Read(reader)
		if err != nil {
			return nil, err
		}

		if checkHeadersEnd(headerBuf) {
			break
		}
		buf = append(buf, headerBuf)
	}
	return buf, nil
}
