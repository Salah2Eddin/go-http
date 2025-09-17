package httpheader

import (
	"bytes"
	"fmt"
	"github.com/Salah2Eddin/go-http/pkg/util/charutil"
	"strings"
)

type Header struct {
	name   string
	values []HeaderValue
}

type HeaderValue struct {
	value  string
	params map[string]string
}

const (
	valueSeparator = ","
	paramSeparator = ";"
)

func needQuotes(s string) bool {
	return strings.ContainsFunc(s, func(r rune) bool { return !charutil.IsTChar(byte(r)) })
}

func quoteString(s string) string {
	// replace any \ with \\ and " with \"
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", s)
}

func formatForHeaders(s string) string {
	if needQuotes(s) {
		return quoteString(s)
	}
	return s
}

func NewHeaderValue(valueByte []byte, paramsBytes []byte) HeaderValue {
	value := string(valueByte)
	params := make(map[string]string)
	for _, param := range bytes.Split(paramsBytes, []byte(";")) {
		k, v, exists := bytes.Cut(param, []byte("="))
		if !exists {
			continue
		}
		params[string(k)] = string(v)
	}
	return HeaderValue{
		value:  value,
		params: params,
	}
}

func (h *HeaderValue) Value() string {
	return h.value
}

func (h *HeaderValue) Params() map[string]string {
	return h.params
}

func (h *HeaderValue) SetParam(name, value string) {
	h.params[name] = value
}

func (h *HeaderValue) GetParam(name string) (string, bool) {
	val, exists := h.params[name]
	return val, exists
}

func (h *HeaderValue) String() string {
	valueStr := formatForHeaders(h.value)
	for key, value := range h.params {
		valueStr += paramSeparator
		valueStr += formatForHeaders(key) + "=" + formatForHeaders(value)
	}
	return valueStr
}

func NewHeader(name string, values []HeaderValue) Header {
	return Header{
		name:   name,
		values: values,
	}
}

func (h *Header) Name() string {
	return h.name
}

func (h *Header) Values() []HeaderValue {
	return h.values
}

func (h *Header) AddValue(value HeaderValue) {
	h.values = append(h.values, value)
}

func (h *Header) String() string {
	headerStr := ""
	for _, headerValue := range h.values {
		if len(headerStr) != 0 {
			headerStr += valueSeparator
		}
		headerStr += headerValue.String()
	}
	return headerStr
}
