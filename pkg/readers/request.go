package readers

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/httpheader"
	"github.com/Salah2Eddin/go-http/pkg/reqline"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"strconv"
)

type RequestReader struct {
	reqLineReader requestLineReader
	headersReader headersReader
}

func (r RequestReader) bodyReaderFactory(headers *httpheader.Headers) iReader[any] {
	//TODO: other body reading strategies
	if val, exists := headers.Get("transfer-encoding"); exists {
		panic("Only supports content-length and empty bodies ")
	} else if val, exists = headers.Get("content-length"); exists {
		length, err := strconv.Atoi(val.Values()[0].Value())
		if err != nil {
			return nil
		}
		return newLengthBodyReader(length)
	} else {
		return EmptyBodyReader{}
	}
}

func (r RequestReader) Parse(reader *bufio.Reader) (request.Request, error) {
	reqLineBuf, err := r.reqLineReader.Read(reader)
	if err != nil {
		return request.Request{}, err
	}
	reqLine, err := reqline.ParseRequestLine(reqLineBuf)
	if err != nil {
		return request.Request{}, err
	}

	headersBuf, err := r.headersReader.Read(reader)
	headers, err := httpheader.ParseRequestHeaders(headersBuf)
	if err != nil {
		return request.Request{}, err
	}

	bodyReader := r.bodyReaderFactory(&headers)
	body, err := bodyReader.Read(reader)
	if err != nil {
		return request.Request{}, err
	}
	// TODO: process body here

	return request.NewRequest(
		reqLine,
		headers,
		&body), nil
}
