package serializers

type IValueFormatter interface {
	format(s string) string
}

var formatters = map[string]IValueFormatter{
	"content-type": IdentityFormatter{},
}

func formatterFactory(key string) IValueFormatter {
	if formatter, ok := formatters[key]; ok {
		return formatter
	}
	return QuotedFormatter{}
}

type QuotedFormatter struct{}

func (QuotedFormatter) format(s string) string {
	if needQuotes(s) {
		return quoteString(s)
	}
	return s
}

type IdentityFormatter struct {
}

func (IdentityFormatter) format(s string) string {
	return s
}
