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

// IsUnreserved checks if a character is an unreserved character in URIs
// Unreserved characters are: ALPHA / DIGIT / "-" / "." / "_" / "~"
func IsUnreserved(b byte) bool {
	switch b {
	case '-', '.', '_', '~':
		return true
	default:
		return IsAlphaNum(b)
	}
}

// IsSubDelim checks if a character is a sub-delimiter in URIs
func IsSubDelim(b byte) bool {
	switch b {
	case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=':
		return true
	default:
		return false
	}
}

// IsGenDelim checks if a character is a gen-delimiter in URIs
func IsGenDelim(b byte) bool {
	switch b {
	case ':', '/', '?', '#', '[', ']', '@':
		return true
	default:
		return false
	}
}

// IsReserved checks if a character is a reserved character in URIs
func IsReserved(b byte) bool {
	return IsGenDelim(b) || IsSubDelim(b)
}

// IsPChar checks if a character is a pchar (path character) in URIs
func IsPChar(b byte) bool {
	return IsUnreserved(b) || IsSubDelim(b) || b == ':' || b == '@'
}

// IsSchemeChar checks if a character is valid in a URI scheme.
func IsSchemeChar(b byte) bool {
	switch b {
	case '+', '-', '.':
		return true
	default:
		return IsAlphaNum(b)
	}
}
