package readers

import "bufio"

func checkCRLF(bytes []byte) bool {
	size := len(bytes)
	if size < 2 {
		return false
	}

	return bytes[size-2] == '\r' && bytes[size-1] == '\n'
}

func readLine(reader *bufio.Reader, buf *[]byte) error {
	for !checkCRLF(*buf) {
		next, err := reader.ReadByte()
		if err != nil {
			return err
		}
		*buf = append(*buf, next)
	}

	// remove CRLF from the line
	*buf = (*buf)[:len(*buf)-2]

	return nil
}

func checkHeadersEnd(bytes []byte) bool {
	return len(bytes) == 0
}
