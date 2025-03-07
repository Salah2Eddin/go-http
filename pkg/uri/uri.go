package uri

import "strings"

type Uri struct {
	uri        string
	path       []string
	parameters map[string]string
}

func parseUriParameters(paramsStr string) map[string]string {
	if paramsStr == "" {
		params := make(map[string]string)
		return params
	}

	paramsParts := strings.Split(paramsStr, "&")
	params := make(map[string]string, len(paramsParts))
	for _, paramKv := range paramsParts {
		key, value, found := strings.Cut(paramKv, "=")

		// ignores invalid uri parameters
		if found {
			params[key] = value
		}
	}
	return params
}

func parseUri(uri string) (string, string) {
	uri, paramsStr, _ := strings.Cut(uri, "?")

	return uri, paramsStr
}

func NewUri(uri string) Uri {
	uri, paramsStr := parseUri(uri)
	params := parseUriParameters(paramsStr)
	path := strings.Split(uri, "/")
	return Uri{
		uri:        uri,
		path:       path,
		parameters: params,
	}
}

func (uri Uri) String() string {
	return uri.uri
}

func (uri Uri) GetParameter(param string) (string, bool) {
	value, found := uri.parameters[param]
	return value, found
}

func (uri Uri) GetPath() []string {
	return uri.path
}
