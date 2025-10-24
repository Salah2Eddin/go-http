package readers

import (
	"bufio"
)

type iReader[T any] interface {
	Read(reader *bufio.Reader) ([]byte, error)
}
