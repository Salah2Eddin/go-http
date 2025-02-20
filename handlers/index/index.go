package index

import (
	"ducky/http/pkg/request"
	"ducky/http/pkg/response"
	statuscodes "ducky/http/pkg/response/statuscodes"
)

func GET(request *request.Request) *response.Response {

	status := statuscodes.Status200()

	headers := response.NewResponseHeaders()
	headers.Set("content-type", "text/plain")

	body := []byte("Hello World!")

	response := response.NewResponse(status, headers, &body)
	return response
}
