package uri

import "strings"

type Uri struct {
	uri        string
	parameters map[string]string
}

func parseUriParameters(params_str string) map[string]string {
	if params_str == "" {
		params := make(map[string]string)
		return params
	}

	params_parts := strings.Split(params_str, "&")
	params := make(map[string]string, len(params_parts))
	for _, param_kv := range params_parts {
		key, value, found := strings.Cut(param_kv, "=")

		// ignores invalid uri parameters
		if found {
			params[key] = value
		}
	}
	return params
}

func parseUri(uri string) (string, string) {
	uri, params_str, _ := strings.Cut(uri, "?")

	return uri, params_str
}

func NewUri(uri string) *Uri {
	uri, params_str := parseUri(uri)
	params := parseUriParameters(params_str)
	return &Uri{
		uri:        uri,
		parameters: params,
	}
}

func (uri *Uri) String() string {
	return uri.uri
}

func (uri *Uri) GetParameter(param string) (string, bool) {
	value, found := uri.parameters[param]
	return value, found
}
