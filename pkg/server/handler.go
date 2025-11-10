package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/readers"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/router"
	"github.com/Salah2Eddin/go-http/pkg/serializers"
	"net"
)

type ConnectionHandler struct {
	conn           net.Conn
	reader         readers.IRequestReader
	serializer     serializers.ISerializer[*response.Response]
	errorResponder IErrorResponder
	router         router.IRouter
}

func (ch *ConnectionHandler) handleError(err *pkgerrors.AppError) *response.Response {
	return ch.errorResponder.From(err)
}

func (ch *ConnectionHandler) handle(req *request.Request) *response.Response {
	handler, routerError := ch.router.GetRequestHandler(req.Uri(), req.Method())
	if routerError != nil {
		return ch.handleError(routerError)
	}
	res, handlerError := handler(req)
	if handlerError != nil {
		return ch.handleError(pkgerrors.NewAppError(handlerError))
	}

	return res
}

func (ch *ConnectionHandler) closeConn() {
	err := ch.conn.Close()
	if err != nil {
		fmt.Printf("Failed to close connection: %s\n", ch.conn.RemoteAddr().String())
	}
}

func (ch *ConnectionHandler) Handle() {
	defer ch.closeConn()

	reader := bufio.NewReader(ch.conn)

	req, parseError := ch.reader.Parse(reader)

	var res *response.Response
	if parseError != nil {
		res = ch.handleError(parseError)
	} else {
		res = ch.handle(req)
	}

	buf := bytes.Buffer{}
	ch.serializer.Serialize(res, &buf)

	_, writeError := ch.conn.Write(buf.Bytes())
	if writeError != nil {
		fmt.Printf("Error writing to conn %s:%s\n", ch.conn.RemoteAddr(), writeError.Error())
	}
}
