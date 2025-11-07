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
	"github.com/Salah2Eddin/go-http/pkg/uri"
	"net"
	"strings"
)

const HTTPVersion = "HTTP/1.0"

type Server struct {
	router         *router.Router
	addr           *Address
	reader         readers.IRequestReader
	serializer     serializers.ISerializer[*response.Response]
	errorResponder IErrorResponder
}

func closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Printf("Failed to close connection: %s\n", conn.RemoteAddr().String())
	}
}

func closeListener(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		panic(err)
	}
}

// NewServer creates and initializes a new Server instance with the provided address or a default address if nil.
func NewServer(address *Address) *Server {
	if address == nil {
		address = &Address{Port: "8576"} // Default address
	}

	// Initialize the server with address and router in one statement
	return &Server{
		addr:           address,
		router:         router.NewRouter(),
		serializer:     serializers.NewResponseSerializer(),
		reader:         readers.NewRequestReader(),
		errorResponder: JSONErrorResponder{},
	}
}

// AddHandler Registers a new handler for the given URI and HTTP method.
// If the route corresponding to the URI does not exist, a new route is created.
func (server *Server) AddHandler(uriStr string, method string, handler router.Handler) error {
	return server.router.AddHandler(uri.NewUri(uriStr), method, handler)
}

func (server *Server) handleError(err *pkgerrors.AppError) *response.Response {
	return server.errorResponder.From(err)
}

func (server *Server) handle(req *request.Request) *response.Response {
	handler, routerError := server.router.GetRequestHandler(req.Uri(), req.Method())
	if routerError != nil {
		return server.handleError(routerError)
	}
	res, handlerError := handler(req)
	if handlerError != nil {
		return server.handleError(pkgerrors.NewAppError(handlerError))
	}

	return res
}

// processConnection Handles an incoming client connection.
// It reads and parses the request, handles it, writes the response to the stream,
// and then closes the connection.
func (server *Server) processConnection(conn net.Conn) {
	defer closeConn(conn)
	reader := bufio.NewReader(conn)

	req, parseError := server.reader.Parse(reader)

	var res *response.Response
	if parseError != nil {
		res = server.handleError(parseError)
	} else {
		res = server.handle(req)
	}

	buf := bytes.Buffer{}
	server.serializer.Serialize(res, &buf)

	_, writeError := conn.Write(buf.Bytes())
	if writeError != nil {
		fmt.Printf("Error writing to conn %s:%s\n", conn.RemoteAddr(), writeError.Error())
	}
}

// Start Initializes the server, listens for incoming connections,
// and handles them concurrently.
func (server *Server) Start() {
	listener, err := net.Listen("tcp4", server.addr.String())
	if err != nil {
		fmt.Println("Error creating listener:", err.Error())
		return
	}
	defer closeListener(listener)

	// update address and port in case they were automatically assigned
	server.addr.IP, server.addr.Port, _ = strings.Cut(listener.Addr().String(), ":")
	fmt.Printf("Listening on: %v\n", server.addr.String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error on connection:", err.Error())
			continue
		}
		fmt.Printf("%s connected\n", conn.RemoteAddr())
		go server.processConnection(conn)
	}
}
