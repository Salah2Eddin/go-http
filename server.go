package main

import (
	"bufio"
	"fmt"
	"net"
)

func handle_request(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		if message == "" {
			break
		}
		fmt.Println(message)
	}
}

func main() {
	listner, err := net.Listen("tcp4", ":3490")
	if err != nil {
		fmt.Println("Error")
	}
	defer listner.Close()

	fmt.Printf("Listening on address: %s\n", listner.Addr())
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Printf("Error %s\n", err.Error())
			continue
		}
		go handle_request(conn)
	}
}
