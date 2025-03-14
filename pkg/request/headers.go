package request

import "strings"

type Headers struct {
	headers map[string]Header
}

func NewRequestHeaders() Headers {
	return Headers{headers: make(map[string]Header)}
}

func (req *Headers) Add(header Header) {
	name := header.Name()
	name = strings.ToLower(name)
	req.headers[name] = header
}

func (req *Headers) Get(name string) (Header, bool) {
	name = strings.ToLower(name)
	val, exists := req.headers[name]
	return val, exists
}
