package httpheader

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/util"
)

func Read(reader *bufio.Reader, buf *[][]byte) error {
	for {
		headerBytes, err := util.ReadLine(reader)
		if err != nil {
			return err
		}

		if checkHeadersEnd(&headerBytes) {
			break
		}
		*buf = append(*buf, headerBytes)
	}
	return nil
}
