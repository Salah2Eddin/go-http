package server

import (
	"bufio"
	"ducky/http/pkg/parsers"
	"ducky/http/pkg/router"
	"fmt"
	"net"
	"strings"
)

type Server struct {
	router *router.Router
	ip     string
	port   string
}

type ServerAddress struct {
	Ip   string
	Port string
}

func NewServer(address *ServerAddress) *Server {
	server := &Server{}
	if address != nil {
		server.ip = address.Ip
		server.port = address.Port
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

func (server *Server) handleRequest(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, err := parsers.ParseRequest(reader)
	if err != nil {
		fmt.Println("Request Error:", err.Error())
		return
	}

	response := server.router.RouteRequest(request)
	conn.Write([]byte(response.String()))
}

func (server *Server) Start() {
	listner, err := net.Listen("tcp4", fmt.Sprintf("%s:%s", server.ip, server.port))
	if err != nil {
		fmt.Println("Error creating listner:", err.Error())
		return
	}
	defer listner.Close()

	// update address and port in case they were automatically assigned
	server.ip, server.port, _ = strings.Cut(listner.Addr().String(), ":")
	fmt.Printf("Listening on: %s:%s\n", server.ip, server.port)

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error on connection:", err.Error())
			continue
		}
		fmt.Printf("%s connected\n", conn.RemoteAddr())
		go server.handleRequest(conn)
	}
}
