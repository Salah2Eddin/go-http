package response

type StatusLine struct {
	Version string
	Code    string
	Phrase  string
}

func NewStatusLine(version, code, phrase string) *StatusLine {
	return &StatusLine{
		Version: version, Code: code, Phrase: phrase,
	}
}
