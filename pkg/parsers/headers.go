package parsers

import (
	"bytes"
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
	"ducky/http/pkg/util"
)

func validHeaderName(name_bytes *[]byte) bool {
	WHITE_SPACE := byte(0x20)

	for _, v := range *name_bytes {
		if !util.IsVisiableUSASCII(v) {
			return false
		}

		// No whitespace is allowed between the field(header) name and colon (RFC9112 5.1)
		if v == WHITE_SPACE {
			return false
		}
	}
	return true
}

func validHeaderValue(value_bytes *[]byte) bool {
	for _, v := range *value_bytes {
		/*
			a recipient of CR, LF, or NUL within a field value
			MUST either reject the message or replace each of those characters with SP.
			Field values containing other CTL characters are also invalid;
			however, recipients MAY retain such characters for
			the sake of robustness when they appear within a safe context
			RFC9110 5.5
		*/
		if util.IsCTLCharASCII(v) {
			return false
		}
	}
	return true
}

func processHeaderName(name_bytes *[]byte) []byte {
	// lowercase to gurantee case insensitivity
	processed_bytes := bytes.ToLower(*name_bytes)
	processed_bytes = bytes.TrimSpace(processed_bytes)

	return processed_bytes
}

func processHeaderValue(value_bytes *[]byte) []byte {
	trimmed_bytes := bytes.TrimSpace(*value_bytes)
	return trimmed_bytes
}

func nameValueSplit(header_line_bytes *[]byte) ([]byte, []byte, bool) {
	// COLON splits header into key and value
	COLON := byte(0x3A)

	return bytes.Cut(*header_line_bytes, []byte{COLON})
}

func parseHeaderLine(b *[]byte) (string, string, error) {
	name_bytes, value_bytes, found := nameValueSplit(b)
	if !found || !validHeaderName(&name_bytes) || !validHeaderValue(&value_bytes) {
		return "", "", errors.ErrInvalidHeader{}
	}

	// header name
	name := string(processHeaderName(&name_bytes))

	// header value
	value := string(processHeaderValue(&value_bytes))

	return name, value, nil
}

func parseRequestHeaders(lines *[]*[]byte) (*request.RequestHeaders, error) {
	headers := request.NewRequestHeaders()
	for _, line := range *lines {
		trimmed_line := bytes.TrimSpace(*line)

		name, value, err := parseHeaderLine(&trimmed_line)
		if err != nil {
			return &request.RequestHeaders{}, err
		}
		headers.Set(name, value)
	}

	return headers, nil
}
