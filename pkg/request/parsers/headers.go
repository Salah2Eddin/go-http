package parsers

import (
	"bufio"
	"bytes"
	"ducky/http/pkg/pkgerrors"
	"ducky/http/pkg/request"
	"ducky/http/pkg/util/charutil"
)

const (
	valueSeparatorByte = byte(0x2C)
	doubleQuotesByte   = byte(0x22)
	escapeByte         = byte(0x5C)
)

func checkHeadersEnd(bytes *[]byte) bool {
	return len(*bytes) == 0
}

func validHeaderName(nameBytes []byte) bool {
	for _, v := range nameBytes {
		if !charutil.IsVisibleASCII(v) {
			return false
		}

		// No whitespace is allowed between the field(header) name and colon (RFC9112 5.1)
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

func processHeaderName(nameBytes []byte) string {
	// lowercase to guarantee case insensitivity
	processedBytes := bytes.ToLower(nameBytes)
	return string(processedBytes)
}

func readUnquotedValue(reader *bytes.Reader) ([]byte, error) {
	value := make([]byte, 0)
	whiteSpaces := make([]byte, 0)

	for reader.Len() > 0 {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == valueSeparatorByte {
			return value, nil
		} else if b == doubleQuotesByte {
			return nil, &pkgerrors.ErrInvalidHeader{}
		} else if charutil.IsWhiteSpaceASCII(b) {
			// leading white spaces are ignored
			if len(value) == 0 {
				continue
			}
			whiteSpaces = append(whiteSpaces, b)
		} else {
			// trailing white spaces are ignored
			value = append(value, whiteSpaces...)
			whiteSpaces = make([]byte, 0)
			value = append(value, b)
		}
	}
	return value, nil
}

func readQuotedValue(reader *bytes.Reader) ([]byte, error) {
	// we already read the first double quotes

	value := make([]byte, 0)
	// look for the matching "
	for reader.Len() > 0 {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		// escaped character is allowed, so just skip the next character
		if b == escapeByte {
			escaped, err := reader.ReadByte()
			value = append(value, b)
			value = append(value, escaped)

			if err != nil {
				return nil, err
			}
		} else if b == doubleQuotesByte {
			return value, nil
		} else {
			value = append(value, b)
		}
	}
	return nil, &pkgerrors.ErrInvalidHeader{}
}

func skipUntilNextListItem(reader *bytes.Reader) error {
	for reader.Len() > 0 {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}
		if b == valueSeparatorByte {
			return nil
		} else if !charutil.IsWhiteSpaceASCII(b) {
			return &pkgerrors.ErrInvalidHeader{}
		}
	}
	return nil
}

func splitHeaderValues(valueBytes []byte) ([][]byte, error) {
	valuesListBytes := make([][]byte, 0)

	valueBytesStream := bytes.NewReader(valueBytes)
	for valueBytesStream.Len() > 0 {
		b, err := valueBytesStream.ReadByte()
		if err != nil {
			return nil, err
		}

		if charutil.IsWhiteSpaceASCII(b) {
			continue
		}

		var value []byte
		if b == doubleQuotesByte {
			value, err = readQuotedValue(valueBytesStream)
			if err != nil {
				return nil, err
			}

			// skip until comma
			err = skipUntilNextListItem(valueBytesStream)
			if err != nil {
				return nil, err
			}
		} else {
			value, err = readUnquotedValue(valueBytesStream)
			if err != nil {
				return nil, err
			}
			valuesListBytes = append(valuesListBytes, value)
		}
		valuesListBytes = append(valuesListBytes, value)
	}

	return valuesListBytes, nil
}

func processHeaderValues(valueBytes []byte) ([]string, error) {
	valuesListBytes, err := splitHeaderValues(valueBytes)
	if err != nil {
		return nil, err
	}

	valuesList := make([]string, 0)
	for _, value := range valuesListBytes {
		trimmed := bytes.TrimSpace(value)
		// Empty elements do not contribute to the count of elements present.
		// RFC9110 5.6.1.2
		if len(trimmed) != 0 {
			valuesList = append(valuesList, string(trimmed))
		}
	}
	/*
		at least one non-empty element is required
		RFC9110 5.6.1.2
	*/
	if len(valuesList) == 0 {
		return nil, &pkgerrors.ErrInvalidHeader{}
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
		return "", []string{}, pkgerrors.ErrInvalidHeader{}
	}

	name := processHeaderName(nameBytes)
	values, err := processHeaderValues(valueBytes)
	if err != nil {
		return "", []string{}, pkgerrors.ErrInvalidHeader{}
	}

	return name, values, nil
}

func getRequestHeaders(reader *bufio.Reader) (request.Headers, error) {
	var headersBytes [][]byte
	for {
		headerBytes, err := readLine(reader)
		if err != nil {
			return request.Headers{}, err
		}

		if checkHeadersEnd(&headerBytes) {
			break
		}
		headersBytes = append(headersBytes, headerBytes)
	}
	return parseRequestHeaders(&headersBytes)
}

func parseRequestHeaders(lines *[][]byte) (request.Headers, error) {
	headers := request.NewRequestHeaders()
	for _, line := range *lines {
		name, value, err := parseHeaderLine(line)
		if err != nil {
			return request.Headers{}, err
		}
		headers.Set(name, value)
	}

	return headers, nil
}
