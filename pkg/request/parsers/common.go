package parsers

import "bufio"

func readLine(reader *bufio.Reader) ([]byte, error) {
	var lineBytes []byte

	for !checkCRLF(&lineBytes) {
		next, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		lineBytes = append(lineBytes, next)
	}

	// remove CLRF from the line
	lineBytes = lineBytes[:len(lineBytes)-2]

	return lineBytes, nil
}

func checkCRLF(bytes *[]byte) bool {
	size := len(*bytes)
	if size < 2 {
		return false
	}

	CL := byte(0x0D)
	RF := byte(0x0A)

	return (*bytes)[size-2] == CL && (*bytes)[size-1] == RF
}
