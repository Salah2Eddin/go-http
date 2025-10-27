package httpheader

import (
	"strings"
)

type Headers struct {
	headers map[string]*Header
}

func New() *Headers {
	return &Headers{headers: make(map[string]*Header)}
}

func (headers *Headers) AddFromString(name string, value string) error {
	header, err := NewHeaderFromString(name, value)
	if err != nil {
		return err
	}
	headers.AddFromHeader(header)
	return nil
}

func (headers *Headers) AddFromHeader(header *Header) {
	name := header.Name()
	name = strings.ToLower(name)
	if h, exists := headers.headers[name]; exists {
		// header with same name exists
		// add new header values to it
		for _, value := range *header.Values() {
			h.AddValue(value)
		}
	}
	headers.headers[name] = header
}

func (headers *Headers) Get(name string) (*Header, bool) {
	name = strings.ToLower(name)
	val, exists := headers.headers[name]
	return val, exists
}

func (headers *Headers) Headers() map[string]*Header {
	return headers.headers
}
