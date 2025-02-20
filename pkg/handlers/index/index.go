package index_handlers

import (
	"ducky/http/pkg/request"
)

func GET(request *request.Request) error {
	println(request.Line.Method)
	return nil
}
