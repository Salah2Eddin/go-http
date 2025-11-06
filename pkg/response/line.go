package response

type StatusLine struct {
	Version string
	Code    int
	Phrase  string
}

func NewStatusLine(version string, code int, phrase string) *StatusLine {
	return &StatusLine{
		Version: version, Code: code, Phrase: phrase,
	}
}
