package parsers

import (
	"bytes"
	"ducky/http/pkg/errors"
	"ducky/http/pkg/request"
	"ducky/http/pkg/util"
)

func validHeaderName(nameBytes []byte) bool {
	for _, v := range nameBytes {
		if !util.IsVisibleASCII(v) {
			return false
		}

		// No whitespace is allowed between the field(header) name and colon (RFC9112 5.1)
		if util.IsWhiteSpaceASCII(v) {
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
		if util.IsCTLCharASCII(v) {
			return false
		}
	}
	return true
}

func processHeaderName(nameBytes []byte) []byte {
	// lowercase to guarantee case insensitivity
	processedBytes := bytes.ToLower(nameBytes)
	processedBytes = bytes.TrimSpace(processedBytes)

	return processedBytes
}

func splitHeaderValues(valueBytes []byte) ([][]byte, error) {
	COMMA := byte(0x2C)
	DQUOTE := byte(0x22)
	BSLASH := byte(0x5C)

	size := len(valueBytes)

	valuesListBytes := make([][]byte, 0)

	for i := 0; i < size; i++ {
		if util.IsWhiteSpaceASCII(valueBytes[i]) {
			continue
		}

		start := i
		if valueBytes[start] == DQUOTE {
			// look for the matching "
			i++
			for i < size && valueBytes[i] != DQUOTE {
				// \" is allowed, so just skip whatever after \
				if valueBytes[i] == BSLASH {
					i++
				}
				i++
			}
			if valueBytes[i] == DQUOTE {
				valuesListBytes = append(valuesListBytes, valueBytes[start:i+1])

				// skip until comma
				for i < size && valueBytes[i] != COMMA {
					i++
				}
			} else {
				// unmatched "
				return nil, &errors.ErrInvalidHeader{}
			}
		} else {
			for i < size && valueBytes[i] != COMMA {
				// not proper quoted string
				if valueBytes[i] == DQUOTE {
					return nil, &errors.ErrInvalidHeader{}
				}
				i++
			}
			valuesListBytes = append(valuesListBytes, valueBytes[start:i])
		}
	}

	return valuesListBytes, nil
}

func processHeaderValues(valueBytes []byte) ([][]byte, error) {
	valuesListBytes, err := splitHeaderValues(valueBytes)
	if err != nil {
		return nil, err
	}

	valuesList := make([][]byte, 0)
	for _, value := range valuesListBytes {
		trimmed := bytes.TrimSpace(value)
		// Empty elements do not contribute to the count of elements present.
		// RFC9110 5.6.1.2
		if len(trimmed) != 0 {
			valuesList = append(valuesListBytes, trimmed)
		}
	}
	/*
		at least one non-empty element is required
		RFC9110 5.6.1.2
	*/
	if len(valuesList) == 0 {
		return nil, &errors.ErrInvalidHeader{}
	}

	return valuesList, nil
}

func nameValueSplit(headerLineBytes []byte) ([]byte, []byte, bool) {
	// COLON splits header into key and value
	COLON := byte(0x3A)

	return bytes.Cut(headerLineBytes, []byte{COLON})
}

func parseHeaderLine(headerLineBytes []byte) (string, []string, error) {
	nameBytes, valueBytes, found := nameValueSplit(headerLineBytes)
	if !found || !validHeaderName(nameBytes) || !validHeaderValue(valueBytes) {
		return "", []string{}, errors.ErrInvalidHeader{}
	}

	// header name
	name := string(processHeaderName(nameBytes))

	// header value
	valuesListBytes, err := processHeaderValues(valueBytes)
	if err != nil {
		return "", []string{}, errors.ErrInvalidHeader{}
	}
	valuesList := make([]string, 0)
	for _, v := range valuesListBytes {
		valuesList = append(valuesList, string(v))
	}

	return name, valuesList, nil
}

func parseRequestHeaders(lines *[][]byte) (request.Headers, error) {
	headers := request.NewRequestHeaders()
	for _, line := range *lines {
		trimmedLine := bytes.TrimSpace(line)

		name, value, err := parseHeaderLine(trimmedLine)
		if err != nil {
			return request.Headers{}, err
		}
		headers.Set(name, value)
	}

	return headers, nil
}
