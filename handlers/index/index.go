package index

import (
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
	statuscodes "github.com/Salah2Eddin/go-http/pkg/response/statuscodes"
)

func GET(request *request.Request) *response.Response {

	status := statuscodes.Status200()

	headers := response.NewResponseHeaders()
	headers.Set("content-type", "text/plain")

	body := []byte("Hello World!")

	response := response.NewResponse(status, headers, &body)
	return response
}
