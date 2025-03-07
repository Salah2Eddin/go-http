package statuscodes

import "ducky/http/pkg/response"

func newStatusLine(version, code, phrase string) response.StatusLine {
	status := response.StatusLine{}
	status.Version = version
	status.Code = code
	status.Phrase = phrase
	return status
}

// TODO: extend with HTTP1.1 Errors
// Dont forget to actually use them

// 2XX status codes

// OK
func Status200() response.StatusLine {
	return newStatusLine("HTTP/1.0", "200", "OK")
}

// Created
func Status201() response.StatusLine {
	return newStatusLine("HTTP/1.0", "201", "Created")
}

// Accepted
func Status202() response.StatusLine {
	return newStatusLine("HTTP/1.0", "202", "Accepted")
}

// No Content
func Status204() response.StatusLine {
	return newStatusLine("HTTP/1.0", "204", "No Content")
}

// 3XX status codes

// Multiple Choices
func Status300() response.StatusLine {
	return newStatusLine("HTTP/1.0", "300", "Multiple Choices")
}

// Moved Permanently
func Status301() response.StatusLine {
	return newStatusLine("HTTP/1.0", "301", "Moved Permanently")
}

// Found
func Status302() response.StatusLine {
	return newStatusLine("HTTP/1.0", "302", "Found")
}

// Not Modified
func Status304() response.StatusLine {
	return newStatusLine("HTTP/1.0", "304", "Not Modified")
}

// 4XX status codes

// Bad Request
func Status400() response.StatusLine {
	return newStatusLine("HTTP/1.0", "400", "Bad Request")
}

// Unauthorized
func Status401() response.StatusLine {
	return newStatusLine("HTTP/1.0", "401", "Unauthorized")
}

// Forbidden
func Status403() response.StatusLine {
	return newStatusLine("HTTP/1.0", "403", "Forbidden")
}

// Not Found
func Status404() response.StatusLine {
	return newStatusLine("HTTP/1.0", "404", "Not Found")
}

// 5XX status codes

// Internal Server Error
func Status500() response.StatusLine {
	return newStatusLine("HTTP/1.0", "500", "Internal Server Error")
}

// Not Implemented
func Status501() response.StatusLine {
	return newStatusLine("HTTP/1.0", "501", "Not Implemented")
}

// Bad Gateway
func Status502() response.StatusLine {
	return newStatusLine("HTTP/1.0", "502", "Bad Gateway")
}

// Service Unavailable
func Status503() response.StatusLine {
	return newStatusLine("HTTP/1.0", "503", "Service Unavailable")
}
