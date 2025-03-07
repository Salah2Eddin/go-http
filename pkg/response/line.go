package response

import "fmt"

type StatusLine struct {
	Version string
	Code    string
	Phrase  string
}

func (status StatusLine) String() string {
	return fmt.Sprintf("%s %s %s\r\n", status.Version, status.Code, status.Phrase)
}

func (status StatusLine) Bytes() []byte {
	return []byte(status.String())
}
