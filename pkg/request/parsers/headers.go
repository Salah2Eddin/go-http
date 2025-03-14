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
	paramSeparatorByte = byte(0x3b)
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
			// , works as a value separator only outside parentheses
			return value, nil
		} else if b == doubleQuotesByte {
			return nil, &pkgerrors.ErrInvalidHeader{}
		} else if charutil.IsWhiteSpaceASCII(b) {
			// leading white spaces are ignored
			if len(value) == 0 {
				continue
			}
			whiteSpaces = append(whiteSpaces, b)
		} else if b == paramSeparatorByte {
			// ; works as a param separator outside parentheses
			break
		} else {
			// trailing white spaces are ignored
			value = append(value, whiteSpaces...)
			whiteSpaces = make([]byte, 0)
			value = append(value, b)
		}
	}
	return value, nil
}

func readParameters(reader *bytes.Reader) ([]byte, error) {
	value := make([]byte, 0)
	// read value parameters
	for reader.Len() > 0 {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if charutil.IsWhiteSpaceASCII(b) {
			// leading spaces before ; are ok
			continue
		} else if b == paramSeparatorByte {
			// read parameter
			for reader.Len() > 0 {
				b, err := reader.ReadByte()
				if err != nil {
					return nil, err
				}
				if b == valueSeparatorByte || b == paramSeparatorByte {
					break
				} else {
					value = append(value, b)
				}
			}
			break
		} else if b == valueSeparatorByte {
			err := reader.UnreadByte()
			if err != nil {
				return nil, err
			}
			break
		} else {
			return nil, &pkgerrors.ErrInvalidHeader{}
		}
	}
	return value, nil
}

func readQuotedValue(reader *bytes.Reader) ([]byte, error) {
	value := make([]byte, 0)
	// number of quotes found
	count := 0
	for reader.Len() > 0 {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		// escaped character is allowed, so just skip the next character
		if b == escapeByte {
			escaped, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			value = append(value, b)
			value = append(value, escaped)
		} else if b == doubleQuotesByte {
			count++
			if count == 2 {
				// quoted values can't be followed by unquoted value or comments
				// example: "quoted"unquoted -> invalid
				for reader.Len() > 0 {
					b, err = peek(reader)
					if err != nil {
						return nil, err
					}
					// skip white spaces
					if charutil.IsWhiteSpaceASCII(b) {
						_, err = reader.ReadByte()
						if err != nil {
							return nil, err
						}
						continue
					}
					// only parameters and other values are allowed after a quoted value
					if b == paramSeparatorByte || b == valueSeparatorByte {
						return value, nil
					} else {
						return nil, &pkgerrors.ErrInvalidHeader{}
					}
				}
				return value, nil
			}
		} else {
			value = append(value, b)
		}
	}
	return nil, &pkgerrors.ErrInvalidHeader{}
}

func readNextValue(reader *bytes.Reader) ([]byte, []byte, error) {
	value := make([]byte, 0)
	for reader.Len() > 0 {
		b, err := peek(reader)
		if err != nil {
			return nil, nil, err
		}
		if charutil.IsWhiteSpaceASCII(b) {
			_, err := reader.ReadByte()
			if err != nil {
				return nil, nil, err
			}
			continue
		}

		if b == doubleQuotesByte {
			value, err = readQuotedValue(reader)
		} else {
			value, err = readUnquotedValue(reader)
		}
		if err != nil {
			return nil, nil, err
		}

		var params []byte
		if reader.Len() > 0 {
			b, err = peek(reader)
			if b == paramSeparatorByte {
				params, err = readParameters(reader)
				if err != nil {
					return nil, nil, err
				}
			}
		}
		return value, params, nil
	}
	return nil, nil, &pkgerrors.ErrInvalidHeader{}
}

func splitHeaderValues(valueBytes []byte) ([][]byte, [][]byte, error) {
	valueBytesStream := bytes.NewReader(valueBytes)
	values := make([][]byte, 0)
	params := make([][]byte, 0)
	for valueBytesStream.Len() > 0 {
		value, param, err := readNextValue(valueBytesStream)
		if err != nil {
			return nil, nil, err
		}
		values = append(values, value)
		params = append(params, param)
	}
	return values, params, nil
}

func processHeaderValues(valueBytes []byte) ([]request.HeaderValue, error) {
	values, params, err := splitHeaderValues(valueBytes)
	if err != nil {
		return nil, err
	}
	headerValues := make([]request.HeaderValue, 0)
	for i := range values {
		// Empty elements do not contribute to the count of elements present.
		// RFC9110 5.6.1.2
		if len(values[i]) == 0 {
			continue
		}
		value := request.NewHeaderValue(values[i], params[i])
		headerValues = append(headerValues, value)
	}
	/*
		at least one non-empty element is required
		RFC9110 5.6.1.2
	*/
	if len(values) == 0 {
		return nil, &pkgerrors.ErrInvalidHeader{}
	}

	return headerValues, nil
}

func nameValueSplit(headerLineBytes []byte) ([]byte, []byte, bool) {
	// COLON splits header into key and value
	COLON := byte(0x3A)

	return bytes.Cut(headerLineBytes, []byte{COLON})
}

func parseHeaderLine(headerLineBytes []byte) (request.Header, error) {
	nameBytes, valueBytes, found := nameValueSplit(headerLineBytes)
	if !found || !validHeaderName(nameBytes) || !validHeaderValue(valueBytes) {
		return request.Header{}, pkgerrors.ErrInvalidHeader{}
	}

	name := processHeaderName(nameBytes)
	values, err := processHeaderValues(valueBytes)
	if err != nil {
		return request.Header{}, pkgerrors.ErrInvalidHeader{}
	}
	header := request.NewHeader(name, values)
	return header, nil
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
		header, err := parseHeaderLine(line)
		if err != nil {
			return request.Headers{}, err
		}
		headers.Add(header)
	}

	return headers, nil
}
