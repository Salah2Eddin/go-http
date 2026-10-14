package server

import (
	"bufio"
	"bytes"
	"github.com/Salah2Eddin/go-http/pkg/router"
	"io"
	"net"
	"time"

	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/uri"
)

// Mock implementations

type mockConn struct {
	readData   []byte
	writeData  bytes.Buffer
	closed     bool
	closeError error
	writeError error
	remoteAddr net.Addr
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if len(m.readData) == 0 {
		return 0, io.EOF
	}
	n = copy(b, m.readData)
	m.readData = m.readData[n:]
	return n, nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	if m.writeError != nil {
		return 0, m.writeError
	}
	return m.writeData.Write(b)
}

func (m *mockConn) Close() error {
	m.closed = true
	return m.closeError
}

func (m *mockConn) LocalAddr() net.Addr                { return m.remoteAddr }
func (m *mockConn) RemoteAddr() net.Addr               { return m.remoteAddr }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

type mockAddr struct{}

func (m mockAddr) Network() string { return "tcp" }
func (m mockAddr) String() string  { return "127.0.0.1:8080" }

type mockRequestReader struct {
	request *request.Request
	err     *pkgerrors.AppError
}

func (m *mockRequestReader) Parse(reader *bufio.Reader) (*request.Request, *pkgerrors.AppError) {
	return m.request, m.err
}

type mockSerializer struct {
	serialized bool
	buffer     *bytes.Buffer
}

func (m *mockSerializer) Serialize(res *response.Response, buf *bytes.Buffer) {
	m.serialized = true
	m.buffer = buf
	buf.WriteString("HTTP/1.1 200 OK\r\n\r\n")
}

type mockErrorResponder struct {
	response *response.Response
	lastErr  *pkgerrors.AppError
}

func (m *mockErrorResponder) From(err *pkgerrors.AppError) *response.Response {
	m.lastErr = err
	return m.response
}

type mockRouter struct {
	handler      func(*request.Request) (*response.Response, error)
	handlerError *pkgerrors.AppError
	lastUri      *uri.Uri
	lastMethod   string
}

func (m *mockRouter) AddHandler(uri *uri.Uri, method string, handler router.Handler) error {
	return nil
}

func (m *mockRouter) GetRequestHandler(uri *uri.Uri, method string) (router.Handler, *pkgerrors.AppError) {
	m.lastUri = uri
	m.lastMethod = method
	return m.handler, m.handlerError
}

type mockListener struct {
	closed     bool
	closeError error
	acceptConn net.Conn
	acceptErr  error
	acceptOnce bool
	addr       net.Addr
}

func (m *mockListener) Accept() (net.Conn, error) {
	if m.acceptOnce {
		// Block forever to prevent infinite loop in tests
		select {}
	}
	m.acceptOnce = true
	return m.acceptConn, m.acceptErr
}

func (m *mockListener) Close() error {
	m.closed = true
	return m.closeError
}

func (m *mockListener) Addr() net.Addr {
	return m.addr
}

type mockListenerAddr struct{}

func (m mockListenerAddr) Network() string { return "tcp" }
func (m mockListenerAddr) String() string  { return "127.0.0.1:8576" }
