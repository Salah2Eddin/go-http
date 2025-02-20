package main

import (
	"bufio"
	"ducky/http/pkg/parsers"
	"fmt"
	"net"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, err := parsers.ParseRequest(reader)
	if err != nil {
		fmt.Println("Request Error:", err.Error())
		return
	}

	println(request.Line.Method)
}

func main() {
	listner, err := net.Listen("tcp4", ":3490")
	if err != nil {
		fmt.Println("Error on listening:", err.Error())
	}
	defer listner.Close()

	fmt.Println("Listening on: ", listner.Addr())
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error on connection acceptance:", err.Error())
			continue
		}
		fmt.Printf("%s Connected\n", conn.RemoteAddr())
		go handleRequest(conn)
	}
}
