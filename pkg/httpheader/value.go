package httpheader

import (
	"bytes"
	"fmt"
	"strings"
)

type Value struct {
	value  string
	params map[string]string
}

func NewHeaderValues(value string) []Value {
	values, err := processHeaderValues([]byte(value))
	if err != nil {
		return []Value{}
	}
	return values
}

func NewHeaderValueFromBytes(valueByte []byte, paramsBytes []byte) Value {
	value := string(valueByte)
	params := make(map[string]string)
	for _, param := range bytes.Split(paramsBytes, []byte(";")) {
		k, v, exists := bytes.Cut(param, []byte("="))
		if !exists {
			continue
		}
		params[string(k)] = string(v)
	}
	return Value{
		value:  value,
		params: params,
	}
}

func (h *Value) Value() string {
	return h.value
}

func (h *Value) Params() map[string]string {
	return h.params
}

func (h *Value) SetParam(name, value string) {
	h.params[name] = value
}

func (h *Value) GetParam(name string) (string, bool) {
	val, exists := h.params[name]
	return val, exists
}

func (h *Value) DeleteParam(name string) {
	delete(h.params, name)
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

func (h *Value) String() string {
	valueStr := formatForHeaders(h.value)
	for key, value := range h.params {
		valueStr += paramSeparator
		valueStr += formatForHeaders(key) + "=" + formatForHeaders(value)
	}
	return valueStr
}
