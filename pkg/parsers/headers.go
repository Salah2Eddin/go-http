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

func splitHeaderValues(value_bytes []byte) ([][]byte, error) {
	COMMA := byte(0x2C)
	DQUOTE := byte(0x22)
	BSLASH := byte(0x5C)

	size := len(value_bytes)

	values_list_bytes := make([][]byte, 0)

	for i := 0; i < size; i++ {
		if util.IsWhiteSpaceASCII(value_bytes[i]) {
			continue
		}

		start := i
		if value_bytes[start] == DQUOTE {
			// look for the matching "
			i++
			for i < size && value_bytes[i] != DQUOTE {
				i++
				// \" is allowed, so just skip whatever after \
				if value_bytes[i] == BSLASH {
					i++
				}
			}
			if value_bytes[i] == DQUOTE {
				// take all inbetween quotes except empty spaces
				values_list_bytes = append(values_list_bytes, value_bytes[start:i+1])

				// skip until comma
				for i < size && value_bytes[i] != COMMA {
					i++
				}
			} else {
				// unmatched "
				return nil, &errors.ErrInvalidHeader{}
			}
		} else {
			for i < size && value_bytes[i] != COMMA {
				i++

				// not proper quoted string
				if value_bytes[i] == DQUOTE {
					return nil, &errors.ErrInvalidHeader{}
				}
			}
			values_list_bytes = append(values_list_bytes, value_bytes[start:i])
		}
	}

	return values_list_bytes, nil
}

func processHeaderValues(value_bytes *[]byte) ([][]byte, error) {
	values_list_bytes, err := splitHeaderValues(*value_bytes)
	if err != nil {
		return nil, err
	}

	values_list := make([][]byte, 0)
	for _, value := range values_list_bytes {
		trimmed := bytes.TrimSpace(value)
		// Empty elements do not contribute to the count of elements present.
		// RFC9110 5.6.1.2
		if len(trimmed) != 0 {
			values_list = append(values_list_bytes, trimmed)
		}
	}
	/*
		at least one non-empty element is required
		RFC9110 5.6.1.2
	*/
	if len(values_list) == 0 {
		return nil, &errors.ErrInvalidHeader{}
	}

	return values_list, nil
}

func nameValueSplit(header_line_bytes *[]byte) ([]byte, []byte, bool) {
	// COLON splits header into key and value
	COLON := byte(0x3A)

	return bytes.Cut(*header_line_bytes, []byte{COLON})
}

func parseHeaderLine(b *[]byte) (string, []string, error) {
	name_bytes, value_bytes, found := nameValueSplit(b)
	if !found || !validHeaderName(&name_bytes) || !validHeaderValue(&value_bytes) {
		return "", []string{}, errors.ErrInvalidHeader{}
	}

	// header name
	name := string(processHeaderName(&name_bytes))

	// header value
	values_list_bytes, err := processHeaderValues(&value_bytes)
	if err != nil {
		return "", []string{}, errors.ErrInvalidHeader{}
	}
	values_list := make([]string, 0)
	for _, v := range values_list_bytes {
		values_list = append(values_list, string(v))
	}

	return name, values_list, nil
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
