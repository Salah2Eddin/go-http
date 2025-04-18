package statuscodes

import "github.com/Salah2Eddin/go-http/pkg/response"

func newStatusLine(version, code, phrase string) *response.StatusLine {
	status := response.StatusLine{}
	status.Version = version
	status.Code = code
	status.Phrase = phrase
	return &status
}

// 2XX status codes

func Status200() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "200", "OK")
}

func Status201() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "201", "Created")
}

func Status202() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "202", "Accepted")
}

func Status204() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "204", "No Content")
}

// 3XX status codes

func Status300() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "300", "Multiple Choices")
}

func Status301() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "301", "Moved Permanently")
}

func Status302() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "302", "Found")
}

func Status304() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "304", "Not Modified")
}

// 4XX status codes

func Status400() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "400", "Bad Request")
}

func Status401() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "401", "Unauthorized")
}

func Status403() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "403", "Forbidden")
}

func Status404() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "404", "Not Found")
}

// 5XX status codes

func Status500() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "500", "Internal Server Error")
}

func Status501() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "501", "Not Implemented")
}

func Status502() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "502", "Bad Gateway")
}

func Status503() *response.StatusLine {
	return newStatusLine("HTTP/1.0", "503", "Service Unavailable")
}
