package reqline

import (
	"bufio"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/uri"
	"github.com/Salah2Eddin/go-http/pkg/util"
	"github.com/Salah2Eddin/go-http/pkg/util/charutil"
	"strings"
)

func validRequestLine(parts []string) bool {
	if len(parts) != 3 {
		return false
	}
	httpVer := parts[2]
	return strings.HasPrefix(httpVer, "HTTP/")
}

func GetRequestLine(reader *bufio.Reader) (RequestLine, error) {
	requestLineBytes, err := util.ReadLine(reader)
	if err != nil {
		return RequestLine{}, err
	}
	return parseRequestLine(requestLineBytes)
}

func parseRequestLine(requestLineBytes []byte) (RequestLine, error) {

	// Request line must contain bytes in the ASCII range only (RFC9112 2.2)
	if !charutil.ValidateAsciiEncoding(requestLineBytes) {
		return RequestLine{}, pkgerrors.ErrInvalidRequestLine{}
	}

	requestLine := string(requestLineBytes)
	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Fields(requestLine)

	if !validRequestLine(parts) {
		return RequestLine{}, pkgerrors.ErrInvalidRequestLine{}
	}

	method := parts[0]
	uriString := parts[1]
	httpVer := parts[2]

	if !uri.ValidateURI(uriString) {
		return RequestLine{}, &pkgerrors.ErrInvalidUri{Uri: uriString}
	}

	uriObj := uri.NewUri(uriString)
	return NewRequestLine(
		method, // method
		uriObj,
		httpVer, // http version
	), nil
}
