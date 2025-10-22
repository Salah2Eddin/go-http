package httpheader

import (
	"fmt"
)

type ISerializer interface {
	Serialize() string
}

func formatterFactory(key string) IValueFormatter {
	var formatters = map[string]IValueFormatter{
		"content-type": NoFormattingFormatter{},
	}
	if formatter, ok := formatters[key]; ok {
		return formatter
	}
	return DefaultFormatter{}
}

func (h *Header) Serialize() string {
	headerStr := fmt.Sprintf("%s: ", h.name)
	for i := range h.values {
		if i != 0 {
			headerStr += valueSeparator
		}
		formatter := formatterFactory(h.name)
		serializer := newValueSerializer(&h.values[i], formatter)
		headerStr += serializer.Serialize()
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
