package parsers

import (
	"bufio"
	"ducky/http/pkg/request"
	"io"
	"strconv"
	"strings"
)

func ParseRequest(reader *bufio.Reader) (*request.Request, error) {
	// request line
	line, err := reader.ReadString('\n')
	if err != nil {
		return &request.Request{}, err
	}
	requst_line, err := ParseRequestLine(line)
	if err != nil {
		return &request.Request{}, err
	}

	// request headers
	var headers []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return &request.Request{}, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		headers = append(headers, line)
	}

	request_headers, err := ParseRequestHeaders(headers)
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
	request := request.NewRequest(requst_line, request_headers, &request_body)
	return request, nil
}
