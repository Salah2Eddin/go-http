package router

import (
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/uri"
)

type IRouter interface {
	AddHandler(uri *uri.Uri, method string, handler Handler) error
	GetRequestHandler(uri *uri.Uri, method string) (Handler, *pkgerrors.AppError)
}
