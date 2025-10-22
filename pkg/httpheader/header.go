package httpheader

import (
	"strings"

	"github.com/Salah2Eddin/go-http/pkg/util/charutil"
)

type Header struct {
	name   string
	values []Value
}

const (
	valueSeparator = ","
	paramSeparator = ";"
)

func needQuotes(s string) bool {
	return strings.ContainsFunc(s, func(r rune) bool { return !charutil.IsTChar(byte(r)) })
}

func NewHeader(name string, values []Value) Header {
	return Header{
		name:   name,
		values: values,
	}
}

func NewHeaderFromString(name string, value string) (Header, error) {
	values, err := processHeaderValues([]byte(value))
	if err != nil {
		return Header{}, err
	}
	return NewHeader(name, values), err
}

func (h *Header) Name() string {
	return h.name
}

func (h *Header) Values() []Value {
	return h.values
}

func (h *Header) AddValue(value Value) {
	h.values = append(h.values, value)
}

func (h *Header) AddValues(values []Value) {
	for _, value := range values {
		h.values = append(h.values, value)
	}
}

func (h *Header) String() string {
	headerStr := ""
	for _, value := range h.values {
		if len(headerStr) != 0 {
			headerStr += valueSeparator
		}
		headerStr += value.String()
	}
	return headerStr
}
