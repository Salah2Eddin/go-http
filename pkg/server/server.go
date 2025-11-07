package server

import (
	"fmt"
	"github.com/Salah2Eddin/go-http/pkg/readers"
	"github.com/Salah2Eddin/go-http/pkg/response"
	"github.com/Salah2Eddin/go-http/pkg/router"
	"github.com/Salah2Eddin/go-http/pkg/serializers"
	"github.com/Salah2Eddin/go-http/pkg/uri"
	"net"
	"strings"
)

type Server struct {
	router         router.IRouter
	addr           *Address
	reader         readers.IRequestReader
	serializer     serializers.ISerializer[*response.Response]
	errorResponder IErrorResponder
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

// acceptLoop continuously accepts incoming network connections and delegates handling to a new goroutine for each connection.
func (server *Server) acceptLoop(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error on connection:", err.Error())
			continue
		}
		fmt.Printf("%s connected\n", conn.RemoteAddr())

		go func() {
			handler := ConnectionHandler{
				conn:           conn,
				reader:         server.reader,
				serializer:     server.serializer,
				errorResponder: server.errorResponder,
				router:         server.router}
			handler.Handle()
		}()
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

	// Start acceptance loop
	server.acceptLoop(listener)
}
