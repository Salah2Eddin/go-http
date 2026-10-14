package serializers

import (
	"fmt"
	"github.com/Salah2Eddin/go-http/pkg/util/charutil"
	"strings"
)

const (
	headerValueSeparator = ","
	headerParamSeparator = ";"
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
