package httpheader

import (
	"fmt"
	"strings"
)

type DefaultValueSerializer struct {
	v *Value
}

func (d DefaultValueSerializer) quoteString(s string) string {
	// replace any \ with \\ and " with \"
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", s)
}

func (d DefaultValueSerializer) formatForHeaders(s string) string {
	if needQuotes(s) {
		return d.quoteString(s)
	}
	return s
}

func (d DefaultValueSerializer) Serialize() string {
	valueStr := d.formatForHeaders(d.v.value)
	for key, value := range d.v.params {
		valueStr += paramSeparator
		valueStr += d.formatForHeaders(key) + "=" + d.formatForHeaders(value)
	}
	return valueStr
}
