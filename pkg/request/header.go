package request

import "bytes"

type Header struct {
	name   string
	values []HeaderValue
}

type HeaderValue struct {
	value  string
	params map[string]string
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
