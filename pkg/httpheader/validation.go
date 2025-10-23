package httpheader

import "github.com/Salah2Eddin/go-http/pkg/util/charutil"

func validHeaderName(nameBytes []byte) bool {
	for _, v := range nameBytes {
		if !charutil.IsVisibleASCII(v) {
			return false
		}

		// No whitespace is allowed between the field(headers) name and colon (RFC9112 5.1)
		if charutil.IsWhiteSpaceASCII(v) {
			return false
		}
	}
	return true
}

func validHeaderValue(valueBytes []byte) bool {
	for _, v := range valueBytes {
		/*
			a recipient of CR, LF, or NUL within a field value
			MUST either reject the message or replace each of those characters with SP.
			Field values containing other CTL characters are also invalid;
			however, recipients MAY retain such characters for
			the sake of robustness when they appear within a safe context
			RFC9110 5.5
		*/
		if charutil.IsCTLCharASCII(v) {
			return false
		}
	}
	return true
}
