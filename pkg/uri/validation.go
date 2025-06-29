package uri

import (
	"github.com/Salah2Eddin/go-http/pkg/util/charutil"
	"strings"
)

// ValidateURI checks if a given string is a valid URI (absolute or relative)
func ValidateURI(uri string) bool {
	// Empty URI is invalid
	if uri == "" {
		return false
	}

	// Check for absolute URI
	if isAbsoluteURI(uri) {
		return true
	}

	// Check for relative URI
	return isRelativeURI(uri)
}

// isAbsoluteURI checks if the given URI is an absolute URI
func isAbsoluteURI(uri string) bool {
	// Find the scheme component
	schemeEndIndex := strings.Index(uri, schemeDelimiter)
	if schemeEndIndex == -1 {
		return false
	}

	// Validate the scheme
	scheme := uri[:schemeEndIndex]
	if !isValidScheme(scheme) {
		return false
	}

	// Check what comes after the scheme
	remainingURI := uri[schemeEndIndex+1:]

	// For absolute URIs, after the scheme should follow "//" and then authority or path
	if !strings.HasPrefix(remainingURI, authorityPrefix) {
		return false
	}

	// Everything after the authority should be a valid path
	pathStartIndex := strings.IndexAny(remainingURI[2:], "/?#")
	var path string

	if pathStartIndex == -1 {
		// No path, query or fragment
		authority := remainingURI[2:]
		return isValidAuthority(authority)
	} else {
		authority := remainingURI[2 : pathStartIndex+2]
		path = remainingURI[pathStartIndex+2:]

		return isValidAuthority(authority) && isValidPath(path)
	}
}

// isRelativeURI checks if the given URI is a valid relative URI
func isRelativeURI(uri string) bool {
	// A relative URI can be a path, optionally followed by ? and query and/or # and fragment
	return isValidPath(uri)
}

// isValidScheme checks if the scheme component is valid
func isValidScheme(scheme string) bool {
	if len(scheme) == 0 {
		return false
	}

	// First character must be a letter
	if !charutil.IsAlpha(scheme[0]) {
		return false
	}

	// Remaining characters must be letters, digits, plus, period, or hyphen
	for i := 1; i < len(scheme); i++ {
		if !charutil.IsSchemeChar(scheme[i]) {
			return false
		}
	}

	return true
}

// isValidAuthority checks if the authority component is valid
func isValidAuthority(authority string) bool {
	if len(authority) == 0 {
		return true // Empty authority is valid
	}

	// Authority may contain userinfo, host, and port
	// userinfo@host:port

	// Check for userinfo
	userInfoEndIndex := strings.Index(authority, userInfoDelimiter)
	var host string

	if userInfoEndIndex != -1 {
		userInfo := authority[:userInfoEndIndex]
		if !isValidUserInfo(userInfo) {
			return false
		}
		host = authority[userInfoEndIndex+1:]
	} else {
		host = authority
	}

	// Check for port
	portStartIndex := strings.LastIndex(host, portDelimiter)
	if portStartIndex != -1 {
		port := host[portStartIndex+1:]
		if !isValidPort(port) {
			return false
		}
		host = host[:portStartIndex]
	}

	// Host validation
	return isValidHost(host)
}

// isValidUserInfo checks if the userinfo component is valid
func isValidUserInfo(userInfo string) bool {
	// Userinfo can be username:password
	// We just check if it contains only allowed characters
	for i := 0; i < len(userInfo); i++ {
		c := userInfo[i]
		if !charutil.IsUserInfoChar(c) {
			return false
		}
	}
	return true
}

// isValidHost checks if the host component is valid
func isValidHost(host string) bool {
	if len(host) == 0 {
		return false
	}

	// Check for IPv6 address
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return isValidIPv6(host[1 : len(host)-1])
	}

	// Check if it's an IPv4 address
	if isValidIPv4(host) {
		return true
	}

	// Check if it's a domain name
	return isValidDomainName(host)
}

// isValidIPv4 checks if the string is a valid IPv4 address
func isValidIPv4(ip string) bool {
	parts := strings.Split(ip, IPV4Delimiter)
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}

		// Check for leading zeros (invalid)
		if len(part) > 1 && part[0] == '0' {
			return false
		}

		// Parse the number
		num := 0
		for i := 0; i < len(part); i++ {
			if !charutil.IsDigit(part[i]) {
				return false
			}
			num = num*10 + int(part[i]-'0')
		}

		if num > 255 {
			return false
		}
	}

	return true
}

// isValidIPv6 performs a simplified check for IPv6 addresses
func isValidIPv6(ip string) bool {
	segments := strings.Split(ip, IPV6Delimiter)

	// IPv6 address must have between 1 and 8 segments
	if len(segments) < 1 || len(segments) > 8 {
		return false
	}

	// Check if there's a :: shorthand
	doubleColonCount := 0
	for i := 0; i < len(segments); i++ {
		if segments[i] == "" {
			doubleColonCount++
		}
	}

	// At most one :: shorthand is allowed
	if doubleColonCount > 1 {
		return false
	}

	// Check each segment
	for _, segment := range segments {
		if len(segment) == 0 {
			continue // This is part of the :: shorthand
		}

		if len(segment) > 4 {
			return false
		}

		for i := 0; i < len(segment); i++ {
			if !charutil.IsHexDigit(segment[i]) {
				return false
			}
		}
	}

	return true
}

// isValidDomainName checks if the string is a valid domain name
func isValidDomainName(domain string) bool {
	parts := strings.Split(domain, DomainDelimiter)

	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false
		}

		for i := 0; i < len(part); i++ {
			c := part[i]
			if !charutil.IsAlphaNum(c) && c != '-' {
				return false
			}
		}

		// First and last character cannot be a hyphen
		if part[0] == '-' || part[len(part)-1] == '-' {
			return false
		}
	}

	return true
}

// isValidPort checks if the port component is valid
func isValidPort(port string) bool {
	if len(port) == 0 {
		return false
	}

	for i := 0; i < len(port); i++ {
		if !charutil.IsDigit(port[i]) {
			return false
		}
	}

	// Convert port to number to check range
	num := 0
	for i := 0; i < len(port); i++ {
		num = num*10 + int(port[i]-'0')
	}

	return num <= portMax && num >= portMin
}

// isValidPath checks if the path component is valid
func isValidPath(path string) bool {
	// Path can have segments separated by "/"
	// Each segment must have valid characters

	// Handle query and fragment if present
	queryIndex := strings.Index(path, queryDelimiter)
	fragmentIndex := strings.Index(path, fragmentDelimiter)

	// Extract the path part
	pathOnly := path
	if queryIndex != -1 {
		if fragmentIndex != -1 && fragmentIndex < queryIndex {
			// Invalid: fragment before query
			return false
		}

		pathOnly = path[:queryIndex]

		// Validate query
		queryPart := ""
		if fragmentIndex != -1 {
			queryPart = path[queryIndex+1 : fragmentIndex]
		} else {
			queryPart = path[queryIndex+1:]
		}

		if !isValidQuery(queryPart) {
			return false
		}
	}

	if fragmentIndex != -1 {
		if queryIndex != -1 && queryIndex > fragmentIndex {
			// Already handled above
		} else {
			pathOnly = path[:fragmentIndex]
		}

		// Validate fragment
		fragmentPart := path[fragmentIndex+1:]
		if !isValidFragment(fragmentPart) {
			return false
		}
	}

	// Now validate the path part
	segments := strings.Split(pathOnly, pathSegmentDelimiter)
	for _, segment := range segments {
		if !isValidSegment(segment) {
			return false
		}
	}

	return true
}

// isValidSegment checks if a path segment is valid
func isValidSegment(segment string) bool {
	// Empty segment is valid (consecutive slashes)
	if len(segment) == 0 {
		return true
	}

	for i := 0; i < len(segment); i++ {
		c := segment[i]

		// Check for percent-encoding
		if c == '%' {
			if i+2 >= len(segment) {
				return false // Incomplete percent-encoding
			}

			// Check the two hex digits
			if !charutil.IsHexDigit(segment[i+1]) || !charutil.IsHexDigit(segment[i+2]) {
				return false
			}

			i += 2 // Skip the two hex digits
			continue
		}

		// Check for other allowed characters
		if !charutil.IsPChar(c) {
			return false
		}
	}

	return true
}

// isValidQuery checks if the query component is valid
func isValidQuery(query string) bool {
	for i := 0; i < len(query); i++ {
		c := query[i]

		// Check for percent-encoding
		if c == '%' {
			if i+2 >= len(query) {
				return false // Incomplete percent-encoding
			}

			// Check the two hex digits
			if !charutil.IsHexDigit(query[i+1]) || !charutil.IsHexDigit(query[i+2]) {
				return false
			}

			i += 2 // Skip the two hex digits
			continue
		}

		// Check for other allowed characters
		if !charutil.IsPChar(c) && c != '/' && c != '?' {
			return false
		}
	}

	return true
}

// isValidFragment checks if the fragment component is valid
func isValidFragment(fragment string) bool {
	// Fragment follows the same rules as query
	return isValidQuery(fragment)
}
