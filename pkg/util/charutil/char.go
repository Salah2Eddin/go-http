package charutil

const (
	lowercaseStart     = 'a'
	lowercaseEnd       = 'z'
	uppercaseStart     = 'A'
	uppercaseEnd       = 'Z'
	digitStart         = '0'
	digitEnd           = '9'
	hexAlphaStart      = 'a'
	hexAlphaEnd        = 'f'
	hexAlphaUpperStart = 'A'
	hexAlphaUpperEnd   = 'F'
)

// IsAlpha checks if a character is an alphabetic character (a-z, A-Z)
func IsAlpha(b byte) bool {
	return isInRange(b, lowercaseStart, lowercaseEnd) || isInRange(b, uppercaseStart, uppercaseEnd)
}

// IsDigit checks if a character is a digit (0-9)
func IsDigit(b byte) bool {
	return isInRange(b, digitStart, digitEnd)
}

// IsAlphaNum checks if a character is alphanumeric (a-z, A-Z, 0-9)
func IsAlphaNum(b byte) bool {
	return IsAlpha(b) || IsDigit(b)
}

// IsHexAlpha checks if a character is a hexadecimal alphabet (a-f, A-F)
func IsHexAlpha(b byte) bool {
	return isInRange(b, hexAlphaStart, hexAlphaEnd) || isInRange(b, hexAlphaUpperStart, hexAlphaUpperEnd)
}

// IsHexDigit checks if a character is a hexadecimal digit (0-9, a-f, A-F)
func IsHexDigit(b byte) bool {
	return IsDigit(b) || IsHexAlpha(b)
}
