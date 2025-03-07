package uri

import "strings"

// parseURI parses a given URI string into its components and updates the provided Uri object.
func parseURI(uri string, u *Uri) {
	// Empty URI
	if uri == "" {
		return
	}

	// Check for scheme (indicating absolute URI)
	schemeEndIndex := strings.Index(uri, schemeDelimiter)
	if schemeEndIndex != -1 {
		u.scheme = uri[:schemeEndIndex]
		u.isAbsolute = true
		uri = uri[schemeEndIndex+1:]

		// Check for authority component
		if strings.HasPrefix(uri, authorityPrefix) {
			uri = uri[2:]

			// Find where authority ends (next '/', '?', or '#')
			authorityEndIndex := strings.IndexAny(uri, "/?#")
			var authority string

			if authorityEndIndex == -1 {
				authority = uri
				uri = ""
			} else {
				authority = uri[:authorityEndIndex]
				uri = uri[authorityEndIndex:]
			}

			// Parse authority (userinfo, host, port)
			parseAuthority(authority, u)
		}
	}
	parsePath(uri, u)
}

// parseAuthority parses the authority component of a URI, including userinfo, host, and port, and updates the given Uri.
func parseAuthority(authority string, u *Uri) {
	// Extract userinfo if present
	userInfoEndIndex := strings.Index(authority, "@")
	var host string

	if userInfoEndIndex != -1 {
		u.userInfo = authority[:userInfoEndIndex]
		host = authority[userInfoEndIndex+1:]
	} else {
		host = authority
	}

	// Extract port if present
	portStartIndex := strings.LastIndex(host, ":")

	// Handle IPv6 addresses correctly
	if portStartIndex != -1 && !strings.Contains(host[portStartIndex:], "]") {
		u.port = host[portStartIndex+1:]
		u.host = host[:portStartIndex]
	} else {
		u.host = host
	}
}

// extractPath extracts the path component of a URI, excluding query parameters and fragments.
func extractPath(uri string) string {
	queryIndex := strings.Index(uri, "?")
	fragmentIndex := strings.Index(uri, "#")
	if queryIndex != -1 {
		return uri[:queryIndex]
	} else if fragmentIndex != -1 {
		return uri[:fragmentIndex]
	}
	return uri
}

// parsePath extracts the path, query, and fragment components from the provided URI and assigns them to the Uri object.
func parsePath(uri string, u *Uri) {
	// Parse each component
	u.path = extractPath(uri)
	u.fragment = extractFragment(uri)

	queryString := extractQuery(uri)
	parseQuery(queryString, u)
}

// extractQuery extracts the query string from a given URI, excluding the '?' and any fragment following '#'.
// Returns an empty string if no query is present.
func extractQuery(uri string) string {
	queryIndex := strings.Index(uri, "?")
	fragmentIndex := strings.Index(uri, "#")
	if queryIndex != -1 {
		if fragmentIndex != -1 && fragmentIndex > queryIndex {
			return uri[queryIndex+1 : fragmentIndex]
		}
		return uri[queryIndex+1:]
	}
	return ""
}

// parseQuery parses a query string into a map of key-value pairs using predefined delimiters.
// It returns an empty map if the input string is empty or contains no valid key-value pairs.
func parseQuery(query string, u *Uri) {
	if len(query) == 0 {
		return
	}

	parts := strings.Split(query, queryKeyDelimiter)
	parameters := make(map[string]string, len(parts))
	for _, entry := range parts {
		key, value, found := strings.Cut(entry, queryKeyValueDelimiter)
		if found {
			parameters[key] = value
		}
	}
	u.query = parameters
}

// extractFragment extracts and returns the fragment portion of the given URI string after the '#' symbol.
// If no fragment is present, it returns an empty string.
func extractFragment(uri string) string {
	fragmentIndex := strings.Index(uri, "#")
	if fragmentIndex != -1 {
		return uri[fragmentIndex+1:]
	}
	return ""
}
