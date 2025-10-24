package httpheaders

import (
	"bytes"

	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/util"
	"github.com/Salah2Eddin/go-http/pkg/util/charutil"
)

const (
	valueSeparatorByte = byte(0x2C)
	doubleQuotesByte   = byte(0x22)
	escapeByte         = byte(0x5C)
	paramSeparatorByte = byte(0x3b)
)

func processHeaderName(nameBytes []byte) string {
	// lowercase to guarantee case insensitivity
	processedBytes := bytes.ToLower(nameBytes)
	return string(processedBytes)
}

func parseUnquotedValue(reader *bytes.Reader) ([]byte, error) {
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
			return nil, pkgerrors.ErrInvalidHeader{Reason: "Unquoted value contains quotes"}
		} else if charutil.IsWhiteSpaceASCII(b) {
			// leading white spaces are ignored
			if len(value) == 0 {
				continue
			}
			whiteSpaces = append(whiteSpaces, b)
		} else if b == paramSeparatorByte {
			// ; works as a param separator outside parentheses
			err := reader.UnreadByte()
			if err != nil {
				return nil, err
			}
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

func parseParameters(reader *bytes.Reader) ([]byte, error) {
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
			return nil, pkgerrors.ErrInvalidHeader{Reason: "parameter parsing invalid syntax"}
		}
	}
	return value, nil
}

func parseQuotedValue(reader *bytes.Reader) ([]byte, error) {
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
					b, err = util.Peek(reader)
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
						return nil, pkgerrors.ErrInvalidHeader{Reason: "Unquoted value after quoted value"}
					}
				}
				return value, nil
			}
		} else {
			value = append(value, b)
		}
	}
	return nil, &pkgerrors.ErrInvalidHeader{Reason: "Quoted value never closed"}
}

func parseNextValue(reader *bytes.Reader) ([]byte, []byte, error) {
	value := make([]byte, 0)
	for reader.Len() > 0 {
		b, err := util.Peek(reader)
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
			value, err = parseQuotedValue(reader)
		} else {
			value, err = parseUnquotedValue(reader)
		}
		if err != nil {
			return nil, nil, err
		}

		var params []byte
		if reader.Len() > 0 {
			b, err = util.Peek(reader)
			if b == paramSeparatorByte {
				params, err = parseParameters(reader)
				if err != nil {
					return nil, nil, err
				}
			}
		}
		return value, params, nil
	}
	return nil, nil, pkgerrors.ErrInvalidHeader{Reason: "Header value parsing fails"}
}

func splitHeaderValues(valueBytes []byte) ([][]byte, [][]byte, error) {
	valueBytesStream := bytes.NewReader(valueBytes)
	values := make([][]byte, 0)
	params := make([][]byte, 0)
	for valueBytesStream.Len() > 0 {
		value, param, err := parseNextValue(valueBytesStream)
		if err != nil {
			return nil, nil, err
		}
		values = append(values, value)
		params = append(params, param)
	}
	return values, params, nil
}

func processHeaderValues(valueBytes []byte) ([]*Value, error) {
	values, params, err := splitHeaderValues(valueBytes)
	if err != nil {
		return nil, err
	}
	headerValues := make([]*Value, 0)
	for i := range values {
		// Empty elements do not contribute to the count of elements present.
		// RFC9110 5.6.1.2
		if len(values[i]) == 0 {
			continue
		}
		value := NewHeaderValueFromBytes(values[i], params[i])
		headerValues = append(headerValues, value)
	}
	/*
		at least one non-empty element is required
		RFC9110 5.6.1.2
	*/
	if len(values) == 0 {
		return nil, pkgerrors.ErrInvalidHeader{Reason: "No non-empty header values"}
	}

	return headerValues, nil
}

func nameValueSplit(headerLineBytes []byte) ([]byte, []byte, bool) {
	// COLON splits headers into key and value
	COLON := byte(0x3A)

	return bytes.Cut(headerLineBytes, []byte{COLON})
}

func parseHeaderLine(headerLineBytes []byte) (*Header, error) {
	nameBytes, valueBytes, found := nameValueSplit(headerLineBytes)
	if !found || !validHeaderName(nameBytes) || !validHeaderValue(valueBytes) {
		return nil, pkgerrors.ErrInvalidHeader{Reason: "Reason"}
	}

	name := processHeaderName(nameBytes)
	values, err := processHeaderValues(valueBytes)
	if err != nil {
		return nil, err
	}
	header := NewHeader(name, values)
	return header, nil
}

func ParseRequestHeaders(lines [][]byte) (*Headers, error) {
	headers := New()
	for _, line := range lines {
		header, err := parseHeaderLine(line)
		if err != nil {
			return nil, err
		}
		headers.AddFromHeader(header)
	}

	return headers, nil
}
