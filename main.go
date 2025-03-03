package main

import (
	"ducky/http/pkg/request"
	"ducky/http/pkg/response"
	"ducky/http/pkg/response/statuscodes"
	"ducky/http/pkg/server"
	"fmt"
)

func index(request *request.Request) *response.Response {
	status := statuscodes.Status200()

	headers := response.NewResponseHeaders()
	headers.Add("content-type", "text/html")

	var body []byte
	if name, exists := request.GetUriParameter("name"); exists {
		body = []byte(fmt.Sprintf("<h1>Hello, %s!</h1>", name))
	} else {
		body = []byte("<h1>Hello, World!</h1>")
	}

	return response.NewResponse(status, headers, &body)
}

func main() {
	server := server.NewServer(&server.ServerAddress{Ip: "127.0.0.1", Port: "8008"})
	server.AddHandler("/", "GET", index)
	server.Start()
}
