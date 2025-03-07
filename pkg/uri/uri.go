package uri

import "strings"

type Uri struct {
	uri        string
	scheme     string
	userInfo   string
	host       string
	port       string
	path       string
	query      map[string]string
	fragment   string
	isAbsolute bool
}

// NewUri creates a Uri instance from the provided full URI string by parsing its components.
func NewUri(fullUri string) Uri {
	uri := Uri{uri: fullUri}
	parseURI(fullUri, &uri)
	return uri
}

// String returns the path component of the Uri as a string.
func (u Uri) String() string {
	uri := u.scheme + "://"
	if u.userInfo != "" {
		uri += u.userInfo + "@"
	}
	uri += u.host
	if u.port != "" {
		uri += ":" + u.port
	}
	uri += u.path
	if len(u.query) > 0 {
		var queryParts []string
		for key, value := range u.query {
			queryParts = append(queryParts, key+"="+value)
		}
		uri += "?" + strings.Join(queryParts, "&")
	}
	if u.fragment != "" {
		uri += "#" + u.fragment
	}
	return uri
}

// GetQueryParameter retrieves the value of a query parameter by its name and indicates if it was found in the URI query.
func (u Uri) GetQueryParameter(param string) (string, bool) {
	value, found := u.query[param]
	return value, found
}

// GetSegments splits the URI's path into its individual segments using the path segment delimiter and returns them as a slice.
func (u Uri) GetSegments() []string {
	return strings.Split(u.path, pathSegmentDelimiter)
}
