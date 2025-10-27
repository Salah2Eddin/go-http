package readers

import (
	"bufio"
)

type iReader interface {
	Read(reader *bufio.Reader) ([]byte, error)
}
