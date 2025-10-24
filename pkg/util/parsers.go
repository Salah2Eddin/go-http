package util

import (
	"bytes"
)

func Peek(reader *bytes.Reader) (byte, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return b, err
	}
	err = reader.UnreadByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}
