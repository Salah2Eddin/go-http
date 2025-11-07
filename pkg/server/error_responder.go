package server

import (
	"fmt"
	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/response"
)

type IErrorResponder interface {
	From(err *pkgerrors.AppError) *response.Response
}

type JSONErrorResponder struct {
}

func (J JSONErrorResponder) generateResponseLine(err *pkgerrors.AppError) *response.StatusLine {
	return response.NewStatusLine(err.HTTPStatusCode())
}

func (J JSONErrorResponder) generateResponseHeaders() *httpheaders.Headers {
	headers := httpheaders.New()
	if err := headers.AddFromString("content-type", "application/json"); err != nil {
		panic(fmt.Sprintf("Failed to set static header: %q", err))
	}
	return headers
}

func (J JSONErrorResponder) generateResponseBody(err *pkgerrors.AppError) []byte {
	return []byte(fmt.Sprintf(`{"error":%q}`, err.Error()))

}

func (J JSONErrorResponder) From(err *pkgerrors.AppError) *response.Response {
	return response.NewResponse(
		J.generateResponseLine(err),
		J.generateResponseHeaders(),
		J.generateResponseBody(err),
	)
}
