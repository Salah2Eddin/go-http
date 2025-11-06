package main

import (
	"fmt"
	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/server"
)

func index(request *request.Request) (*response.Response, pkgerrors.HTTPError) {
	status := response.NewStatusLine(server.HTTPVersion, 200, "OK")

	id := request.Uri().GetSegments()[2]

	headers := httpheaders.New()
	err := headers.AddFromString("content-type", "text/html")
    if err != nil {
        return nil, pkgerrors.NewAppError(err)
    }

	var body []byte
	if name, exists := request.GetUriParameter("name"); exists {
		body = []byte(fmt.Sprintf("<h1>Hello, %s!</h1>", name))
	} else {
		body = []byte("<h1>Hello, World!</h1>")
	}

	body = append(body, []byte(fmt.Sprintf("<h1>Your ID is %s</h1>", id))...)
	resp := response.NewResponse(status, headers, body)
	return resp, nil
}

func main() {
	app := server.NewServer(&server.Address{IP: "127.0.0.1", Port: "8008"})
	err := app.AddHandler("/id/*", "GET", index)
	if err != nil {
		panic(err)
	}

	app.Start()
}
