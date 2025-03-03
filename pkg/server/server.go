package server

import (
	"bufio"
	"ducky/http/pkg/errors"
	"ducky/http/pkg/parsers"
	"ducky/http/pkg/response"
	"ducky/http/pkg/response/statuscodes"
	"ducky/http/pkg/router"
	"fmt"
	"net"
	"strings"
)

type Server struct {
	router *router.Router
	addr   *ServerAddress
}

// Creates and initializes a new Server instance.
// If no address is provided, it defaults to an empty ServerAddress.
func NewServer(address *ServerAddress) *Server {
	server := &Server{addr: address}
	if address == nil {
		server.addr = &ServerAddress{}
	}

	server.router = router.NewRouter()

	return server
}

// Registers a new handler for the given URI and HTTP method.
// If the route corresponding to the URI does not exist, a new route is created.
func (server *Server) AddHandler(uri string, method string, handler router.Handler) {
	route, exists := server.router.GetRoute(uri)
	if !exists {
		route = server.router.NewRoute(uri)
	}
	route.AddHandler(method, handler)
}

// Returns the appropriate HTTP status code
// based on the type of error encountered.
func getErrorStatusCode(err error) *response.StatusLine {
	var status_line *response.StatusLine

	switch err.(type) {
	case errors.ErrInvalidHeader, errors.ErrInvalidRequestLine:
		status_line = statuscodes.Status400()
	default:
		status_line = statuscodes.Status500()
	}
	return status_line
}

// Handles an incoming client connection.
// It reads and parses the request, processes it, writes the response,
// and then closes the connection.
func (server *Server) processConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, err := parsers.ParseRequest(reader)
	var res *response.Response
	if err != nil {
		res = response.NewEmptyResponse(getErrorStatusCode(err))
	} else {
		res = server.router.RouteRequest(request)
	}

	conn.Write(res.Bytes())
}

// Initializes the server, listens for incoming connections,
// and handles them concurrently.
func (server *Server) Start() {
	listner, err := net.Listen("tcp4", server.addr.String())
	if err != nil {
		fmt.Println("Error creating listner:", err.Error())
		return
	}
	defer listner.Close()

	// update address and port in case they were automatically assigned
	server.addr.Ip, server.addr.Port, _ = strings.Cut(listner.Addr().String(), ":")
	fmt.Printf("Listening on: %v\n", server.addr)

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error on connection:", err.Error())
			continue
		}
		fmt.Printf("%s connected\n", conn.RemoteAddr())
		go server.processConnection(conn)
	}
}
