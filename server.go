package main

import (
	"bufio"
	index_handlers "ducky/http/pkg/handlers/index"
	"ducky/http/pkg/parsers"
	"ducky/http/pkg/router"
	"fmt"
	"net"
)

func handleRequest(conn net.Conn, router *router.Router) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, err := parsers.ParseRequest(reader)
	if err != nil {
		fmt.Println("Request Error:", err.Error())
		return
	}

	response := router.Route(request)
	conn.Write([]byte(response.String()))
}

func main() {
	listner, err := net.Listen("tcp4", ":3490")
	if err != nil {
		fmt.Println("Error on listening:", err.Error())
	}
	defer listner.Close()

	fmt.Println("Listening on: ", listner.Addr())

	router := router.NewRouter()
	index := router.NewRoute("/")
	index.AddHandler("GET", index_handlers.GET)

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error on connection acceptance:", err.Error())
			continue
		}
		fmt.Printf("%s Connected\n", conn.RemoteAddr())
		go handleRequest(conn, router)
	}
}
