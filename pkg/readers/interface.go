package readers

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/request"
)

type iReader interface {
	Read(reader *bufio.Reader) ([]byte, error)
}

type IRequestReader interface {
	Parse(reader *bufio.Reader) (*request.Request, *pkgerrors.AppError)
}
