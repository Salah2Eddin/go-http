package parsers

import (
	"bufio"
	"ducky/http/pkg/request"
	"io"
	"strconv"
)

func checkCRLF(bytes []byte) bool {
	size := len(bytes)
	if size < 2 {
		return false
	}

	CL := byte(0x0D)
	RF := byte(0x0A)

	return (bytes)[size-2] == CL && (bytes)[size-1] == RF
}

func checkHeadersEnd(bytes []byte) bool {
	return len(bytes) == 2 && checkCRLF(bytes)
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	var lineBytes []byte

	for !checkCRLF(lineBytes) {
		next, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		lineBytes = append(lineBytes, next)
	}

	return lineBytes, nil
}

func ParseRequest(reader *bufio.Reader) (request.Request, error) {
	// request line
	requestLineBytes, err := readLine(reader)
	if err != nil {
		return request.Request{}, err
	}
	requestLine, err := parseRequestLine(requestLineBytes)
	if err != nil {
		return request.Request{}, err
	}

	// request headers_bytes
	var headersBytes [][]byte
	for {
		headerBytes, err := readLine(reader)
		if err != nil {
			return request.Request{}, err
		}

		// TODO: add line folding support

		if checkHeadersEnd(headerBytes) {
			break
		}
		headersBytes = append(headersBytes, headerBytes)
	}

	requestHeaders, err := parseRequestHeaders(&headersBytes)
	if err != nil {
		return request.Request{}, err
	}

	// request body
	var requestBody []byte
	if lengthHeader, exists := requestHeaders.Get("content-length"); exists {
		lengthStr := lengthHeader[0]
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return request.Request{}, err
		}

		// read request body as byte array.
		// interpreting it is left to the HTTP request handler
		requestBody = make([]byte, length)
		_, err = io.ReadFull(reader, requestBody)
		if err != nil {
			return request.Request{}, err
		}
	}
	req := request.NewRequest(requestLine, requestHeaders, &requestBody)
	return req, nil
}
