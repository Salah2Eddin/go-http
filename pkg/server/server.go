package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/readers"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/response/statuscodes"
	"github.com/Salah2Eddin/go-http/pkg/router"
	"github.com/Salah2Eddin/go-http/pkg/serializers"
	"github.com/Salah2Eddin/go-http/pkg/uri"
	"net"
	"strings"
)

type Server struct {
	router     router.Router
	addr       Address
	reader     readers.RequestReader
	serializer serializers.ResponseSerializer
}

// NewServer creates and initializes a new Server instance with the provided address or a default address if nil.
func NewServer(address *Address) Server {
	if address == nil {
		address = &Address{} // Default address
	}

	// Initialize the server with address and router in one statement
	return Server{
		addr:       *address,
		router:     router.NewRouter(),
		serializer: serializers.NewResponseSerializer(),
	}
}

// AddHandler Registers a new handler for the given URI and HTTP method.
// If the route corresponding to the URI does not exist, a new route is created.
func (server *Server) AddHandler(uriStr string, method string, handler router.Handler) {
	err := server.router.AddHandler(uri.NewUri(uriStr), method, handler)
	if err != nil {
		panic(err)
	}
}

// Returns the appropriate HTTP status code
// based on the type of error encountered.
func mapErrorToStatusCode(err error) response.StatusLine {
	switch err.(type) {
	// ErrExpectedEmptyBody indicates that a request body was received when none was expected, which is a client-side error.
	case pkgerrors.ErrInvalidHeader, pkgerrors.ErrInvalidRequestLine, pkgerrors.ErrExpectedEmptyBody:
		return statuscodes.Status400()
	default:
		return statuscodes.Status500()
	}
}

func closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		panic(err)
	}
}

func closeListener(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		panic(err)
	}
}

// Handles an incoming client connection.
// It reads and parses the request, processes it, writes the response,
// and then closes the connection.
func (server *Server) processConnection(conn net.Conn) {
	defer closeConn(conn)
	reader := bufio.NewReader(conn)

	req, err := server.reader.Parse(reader)

	// TODO: Do something better here
	var res *response.Response
	if err != nil {
		res = response.NewEmptyResponse(mapErrorToStatusCode(err))
	} else {
		handler, err := server.router.GetRequestHandler(req.Uri(), req.Method())
		if err != nil {
			res = response.NewEmptyResponse(mapErrorToStatusCode(err))
		}
		res, err = handler(req)
		if err != nil {
			res = response.NewEmptyResponse(mapErrorToStatusCode(err))
		}
	}

	buf := bytes.Buffer{}
	server.serializer.Serialize(res, &buf)

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Error writing to conn %s:%s\n", conn.RemoteAddr(), err.Error())
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
