package charutil

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
