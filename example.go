package main

import (
	"fmt"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/response/statuscodes"
	"github.com/Salah2Eddin/go-http/pkg/server"
)

func index(request *request.Request) *response.Response {
	status := statuscodes.Status200()

	id := request.Uri().GetSegments()[2]

	headers := response.NewResponseHeaders()
	headers.Add("content-type", "text/html")

	var body []byte
	if name, exists := request.GetUriParameter("name"); exists {
		body = []byte(fmt.Sprintf("<h1>Hello, %s!</h1>", name))
	} else {
		body = []byte("<h1>Hello, World!</h1>")
	}

	body = append(body, []byte(fmt.Sprintf("<h1>Your ID is %s</h1>", id))...)

	return response.NewResponse(status, headers, &body)
}

func main() {
	app := server.NewServer(&server.Address{IP: "127.0.0.1", Port: "8008"})
	app.AddHandler("/id/*", "GET", index)
	app.Start()
}
