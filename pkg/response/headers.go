package response

import "fmt"

type ResponseHeaders struct {
	headers map[string]string
}

func NewResponseHeaders() *ResponseHeaders {
	headers := ResponseHeaders{}
	headers.headers = make(map[string]string)
	return &headers
}

func (h *ResponseHeaders) Set(name string, value string) {
	h.headers[name] = value
}

func (h *ResponseHeaders) Get(name string) (string, bool) {
	val, exists := h.headers[name]
	return val, exists
}

func (headers ResponseHeaders) String() string {
	headers_str := ""
	for k, v := range headers.headers {
		header := fmt.Sprintf("%s:%s\r\n", k, v)
		headers_str += header
	}
	return headers_str
}
