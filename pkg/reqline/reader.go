package reqline

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/util"
)

func Read(reader *bufio.Reader, buf *[]byte) error {
	line, err := util.ReadLine(reader)
	if err != nil {
		return err
	}

	buf = &line
	return nil
}
