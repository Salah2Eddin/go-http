package parsers

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/request"
)

func ParseRequest(reader *bufio.Reader) (request.Request, error) {
	requestLine, err := getRequestLine(reader)
	if err != nil {
		return request.Request{}, err
	}

	requestHeaders, err := getRequestHeaders(reader)
	if err != nil {
		return request.Request{}, err
	}

	requestBody, err := getRequestBody(reader, requestHeaders)
	if err != nil {
		return request.Request{}, err
	}
	req := request.NewRequest(requestLine, requestHeaders, requestBody)
	return req, nil
}
