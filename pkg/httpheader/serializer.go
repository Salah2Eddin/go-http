package httpheader

import (
	"fmt"
)

type ISerializer interface {
	Serialize() string
}

func serializerFactory(key string, value *Value) ISerializer {
	var serializers = map[string]func(v *Value) ISerializer{
		// "content-type": func(v *Value) ISerializer { return nil },
	}
	if factory, ok := serializers[key]; ok {
		return factory(value)
	}
	return DefaultValueSerializer{value}
}

func (h *Header) Serialize() string {
	headerStr := fmt.Sprintf("%s: ", h.name)
	for i := range h.values {
		if i != 0 {
			headerStr += valueSeparator
		}
		headerStr += serializerFactory(h.name, &h.values[i]).Serialize()
	}
	return headerStr
}

func (headers *Headers) Serialize() string {
	headersStr := ""
	for _, header := range headers.headers {
		headersStr += fmt.Sprintf("%s\r\n", header.Serialize())
	}
	return headersStr
}
