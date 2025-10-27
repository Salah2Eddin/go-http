package httpheader

import (
	"bytes"
)

type Value struct {
	value  string
	params map[string]string
}

func NewHeaderValues(value string) *[]Value {
	values, err := processHeaderValues([]byte(value))
	if err != nil {
		return nil
	}
	return values
}

func NewHeaderValueFromBytes(valueByte []byte, paramsBytes []byte) *Value {
	value := string(valueByte)
	params := make(map[string]string)
	for _, param := range bytes.Split(paramsBytes, []byte(";")) {
		k, v, exists := bytes.Cut(param, []byte("="))
		if !exists {
			continue
		}
		params[string(k)] = string(v)
	}
	return &Value{
		value:  value,
		params: params,
	}
}

func (v *Value) Value() string {
	return v.value
}

func (v *Value) Params() map[string]string {
	return v.params
}

func (v *Value) SetParam(name, value string) {
	v.params[name] = value
}

func (v *Value) GetParam(name string) (string, bool) {
	val, exists := v.params[name]
	return val, exists
}

func (v *Value) DeleteParam(name string) {
	delete(v.params, name)
}
