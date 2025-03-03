package response

import (
	"fmt"
	"strings"
)

type ResponseHeaders struct {
	headers map[string][]string
}

func NewResponseHeaders() *ResponseHeaders {
	headers := ResponseHeaders{}
	headers.headers = make(map[string][]string)
	return &headers
}

func (h *ResponseHeaders) Add(name string, value string) {
	if _, exists := h.headers[name]; !exists {
		h.headers[name] = make([]string, 0)
	}
	h.headers[name] = append(h.headers[name], value)
}

func (h *ResponseHeaders) Get(name string) ([]string, bool) {
	val, exists := h.headers[name]
	return val, exists
}

func (headers ResponseHeaders) String() string {
	headers_str := ""
	for key, values := range headers.headers {
		header := fmt.Sprintf("%s:", key)
		for i, value := range values {
			if strings.ContainsRune(value, ',') {
				// replace any " with \"
				quoted := strings.ReplaceAll(value, "\"", "\\\"")
				header += fmt.Sprintf("\"%s\"", quoted)
			} else {
				header += value
			}
			if i != len(values)-1 {
				header += ","
			}
		}
		header += "\r\n"
		headers_str += header
	}
	return headers_str
}
