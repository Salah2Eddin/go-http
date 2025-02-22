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

// Creates a new Server object
func NewServer(address *ServerAddress) *Server {
	server := &Server{addr: address}
	if address == nil {
		server.addr = &ServerAddress{}
	}

	server.router = router.NewRouter()

	return server
}

func (server *Server) AddHandler(uri string, method string, handler router.Handler) {
	route, exists := server.router.GetRoute(uri)
	if !exists {
		route = server.router.NewRoute(uri)
	}
	route.AddHandler(method, handler)
}

func getErrorStatusCode(err error) *response.StatusLine {
	var status_line *response.StatusLine

	// gets the correct status code based on the error
	switch err.(type) {
	case errors.ErrInvalidHeader, errors.ErrInvalidRequestLine:
		status_line = statuscodes.Status400()
	default:
		status_line = statuscodes.Status500()
	}
	return status_line
}

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

	conn.Write([]byte(res.String()))
}

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
