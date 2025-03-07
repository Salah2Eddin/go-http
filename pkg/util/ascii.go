package util

const (
	AsciiMin        = 0x00
	AsciiMax        = 0x7F
	AsciiVisibleMin = 0x21
	AsciiVisibleMax = 0x7E
	AsciiCtlMin     = 0x00
	AsciiCtlMax     = 0x1F
	AsciiSpace      = 0x20
	AsciiTab        = 0x09
	AsciiDelete     = 0x7F
)

// isInRange checks whether the given byte falls between the specified minimum and maximum (inclusive).
func isInRange(b, min, max byte) bool {
	return b >= min && b <= max
}

// IsVisibleASCII determines if a given byte represents a visible ASCII character.
func IsVisibleASCII(b byte) bool {
	return isInRange(b, AsciiVisibleMin, AsciiVisibleMax)
}

// IsASCII determines whether a given byte falls within the ASCII range.
func IsASCII(b byte) bool {
	return isInRange(b, AsciiMin, AsciiMax)
}

// IsCTLCharASCII determines if a given ASCII byte is a control character.
func IsCTLCharASCII(b byte) bool {
	return isInRange(b, AsciiCtlMin, AsciiCtlMax) || b == AsciiDelete
}

// IsWhiteSpaceASCII checks if the given byte represents an ASCII space or tab character.
func IsWhiteSpaceASCII(b byte) bool {
	return b == AsciiSpace || b == AsciiTab
}
