package readers

import "bufio"

type headerReader struct {
}

func (h headerReader) Read(reader *bufio.Reader) ([]byte, error) {
	buf := make([]byte, 0)
	err := readLine(reader, &buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
