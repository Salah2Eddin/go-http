package parsers

import (
	"bufio"
	"ducky/http/pkg/request"
	"io"
	"strconv"
)

func checkCRLF(bytes *[]byte) bool {
	size := len(*bytes)
	if size < 2 {
		return false
	}

	CL := byte(0x0D)
	RF := byte(0x0A)

	return (*bytes)[size-2] == CL && (*bytes)[size-1] == RF
}

func checkHeadersEnd(bytes *[]byte) bool {
	return len(*bytes) == 2 && checkCRLF(bytes)
}

func readLine(reader *bufio.Reader) (*[]byte, error) {
	var line_bytes []byte

	for checkCRLF(&line_bytes) {
		next, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		line_bytes = append(line_bytes, next)
	}

	return &line_bytes, nil
}

func ParseRequest(reader *bufio.Reader) (*request.Request, error) {
	// request line
	request_line_bytes, err := readLine(reader)
	if err != nil {
		return &request.Request{}, err
	}
	request_line, err := parseRequestLine(request_line_bytes)
	if err != nil {
		return &request.Request{}, err
	}

	// request headers_bytes
	var headers_bytes []*[]byte
	for {
		header_bytes, err := readLine(reader)
		if err != nil {
			return &request.Request{}, err
		}

		// TODO: add line folding support

		if checkHeadersEnd(header_bytes) {
			break
		}
		headers_bytes = append(headers_bytes, header_bytes)
	}

	request_headers, err := parseRequestHeaders(&headers_bytes)
	if err != nil {
		return &request.Request{}, err
	}

	// request body
	var request_body []byte
	if length_str, exists := request_headers.Get("content-length"); exists {
		length, err := strconv.Atoi(length_str)
		if err != nil {
			return &request.Request{}, err
		}

		// read request body as byte array.
		// interpreting it is left to the HTTP request handler
		request_body = make([]byte, length)
		_, err = io.ReadFull(reader, request_body)
		if err != nil {
			return &request.Request{}, err
		}
	}
	request := request.NewRequest(request_line, request_headers, &request_body)
	return request, nil
}
