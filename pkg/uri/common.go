package uri

const (
	pathSegmentDelimiter   = "/"
	queryKeyDelimiter      = "&"
	queryKeyValueDelimiter = "="
	queryDelimiter         = "?"
	fragmentDelimiter      = "#"
	userInfoDelimiter      = "@"
	DomainDelimiter        = "."
	IPV4Delimiter          = "."
	IPV6Delimiter          = ":"
	schemeDelimiter        = ":"
	portDelimiter          = ":"

	authorityPrefix = "//"

	portMin = 0
	portMax = 65535
)
