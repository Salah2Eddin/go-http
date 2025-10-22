package httpheader

import (
	"fmt"
	"strings"
)

type IValueFormatter interface {
	format(s string) string
}

func quoteString(s string) string {
	// replace any \ with \\ and " with \"
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", s)
}

type ValueSerializer struct {
	value     *Value
	formatter IValueFormatter
}

func newValueSerializer(value *Value, formatter IValueFormatter) ValueSerializer {
	return ValueSerializer{value, formatter}
}

func (serializer ValueSerializer) Serialize() string {
	valueStr := serializer.formatter.format(serializer.value.value)
	for key, value := range serializer.value.params {
		valueStr += paramSeparator
		valueStr += serializer.formatter.format(key) + "=" + serializer.formatter.format(value)
	}
	return valueStr
}

type DefaultFormatter struct{}

func (DefaultFormatter) format(s string) string {
	if needQuotes(s) {
		return quoteString(s)
	}
	return s
}

type NoFormattingFormatter struct {
}

func (NoFormattingFormatter) format(s string) string {
	return s
}
