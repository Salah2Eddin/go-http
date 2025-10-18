package request

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/httpbody"
	"github.com/Salah2Eddin/go-http/pkg/httpheader"
	"github.com/Salah2Eddin/go-http/pkg/reqline"
)

func ParseRequest(reader *bufio.Reader) (Request, error) {
	requestLine, err := reqline.GetRequestLine(reader)
	if err != nil {
		return Request{}, err
	}

	requestHeaders, err := httpheader.GetRequestHeaders(reader)
	if err != nil {
		return Request{}, err
	}

	requestBody, err := httpbody.GetRequestBody(reader, requestHeaders)
	if err != nil {
		return Request{}, err
	}
	req := NewRequest(requestLine, requestHeaders, requestBody)
	return req, nil
}
