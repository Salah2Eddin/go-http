package response

import (
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string][]string
}

func NewResponseHeaders() Headers {
	headers := Headers{}
	headers.headers = make(map[string][]string)
	return headers
}

func (headers *Headers) Add(name string, value string) {
	if _, exists := headers.headers[name]; !exists {
		headers.headers[name] = make([]string, 0)
	}
	headers.headers[name] = append(headers.headers[name], value)
}

func (headers *Headers) Get(name string) ([]string, bool) {
	val, exists := headers.headers[name]
	return val, exists
}

func (headers *Headers) String() string {
	headersStr := ""
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
		headersStr += header
	}
	return headersStr
}

func (headers *Headers) Bytes() []byte {
	return []byte(headers.String())
}
