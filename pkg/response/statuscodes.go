package response

// 2XX status codes

// Status200 OK
func Status200() *StatusLine {
	return newStatusLine("HTTP/1.0", "200", "OK")
}

// Status201 Created
func Status201() *StatusLine {
	return newStatusLine("HTTP/1.0", "201", "Created")
}

// Status202 Accepted
func Status202() *StatusLine {
	return newStatusLine("HTTP/1.0", "202", "Accepted")
}

// Status204 No Content
func Status204() *StatusLine {
	return newStatusLine("HTTP/1.0", "204", "No Content")
}

// 3XX status codes

// Status300 Multiple Choices
func Status300() *StatusLine {
	return newStatusLine("HTTP/1.0", "300", "Multiple Choices")
}

// Status301 Moved Permanently
func Status301() *StatusLine {
	return newStatusLine("HTTP/1.0", "301", "Moved Permanently")
}

// Status302 Found
func Status302() *StatusLine {
	return newStatusLine("HTTP/1.0", "302", "Found")
}

// Status304 Not Modified
func Status304() *StatusLine {
	return newStatusLine("HTTP/1.0", "304", "Not Modified")
}

// 4XX status codes

// Status400 Bad Request
func Status400() *StatusLine {
	return newStatusLine("HTTP/1.0", "400", "Bad Request")
}

// Status401 Unauthorized
func Status401() *StatusLine {
	return newStatusLine("HTTP/1.0", "401", "Unauthorized")
}

// Status403 Forbidden
func Status403() *StatusLine {
	return newStatusLine("HTTP/1.0", "403", "Forbidden")
}

// Status404 Not Found
func Status404() *StatusLine {
	return newStatusLine("HTTP/1.0", "404", "Not Found")
}

// 5XX status codes

// Status500 Internal Server Error
func Status500() *StatusLine {
	return newStatusLine("HTTP/1.0", "500", "Internal Server Error")
}

// Status501 Not Implemented
func Status501() *StatusLine {
	return newStatusLine("HTTP/1.0", "501", "Not Implemented")
}

// Status502 Bad Gateway
func Status502() *StatusLine {
	return newStatusLine("HTTP/1.0", "502", "Bad Gateway")
}

// Status503 Service Unavailable
func Status503() *StatusLine {
	return newStatusLine("HTTP/1.0", "503", "Service Unavailable")
}
