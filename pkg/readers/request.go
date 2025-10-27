package readers

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/reqline"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"strconv"
)

type RequestReader struct {
	reqLineReader requestLineReader
	headersReader headersReader
}

func NewRequestReader() RequestReader {
	return RequestReader{
		reqLineReader: requestLineReader{},
		headersReader: headersReader{},
	}
}

func (r RequestReader) bodyReaderFactory(headers *httpheaders.Headers) (iReader, error) {
	//TODO: other body reading strategies
	if val, exists := headers.Get("transfer-encoding"); exists {
		return nil, pkgerrors.ErrUnsupportedBodyTransferEncoding{}
	} else if val, exists = headers.Get("content-length"); exists {
		length, err := strconv.Atoi((val.Values())[0].Value())
		if err != nil {
			return nil, pkgerrors.ErrInvalidContentLength{}
		}
		return newLengthBodyReader(length), nil
	} else {
		return EmptyBodyReader{}, nil
	}
}

func (r RequestReader) Parse(reader *bufio.Reader) (*request.Request, error) {
	reqLineBuf, err := r.reqLineReader.Read(reader)
	if err != nil {
		return nil, err
	}
	reqLine, err := reqline.ParseRequestLine(reqLineBuf)
	if err != nil {
		return nil, err
	}

	headersBuf, err := r.headersReader.Read(reader)
	headers, err := httpheaders.ParseRequestHeaders(headersBuf)
	if err != nil {
		return nil, err
	}

	bodyReader, err := r.bodyReaderFactory(headers)
	if err != nil {
		return nil, err
	}
	body, err := bodyReader.Read(reader)
	if err != nil {
		return nil, err
	}
	// TODO: process body here

	req := request.NewRequest(
		reqLine,
		headers,
		&body)
	return &req, nil
}
