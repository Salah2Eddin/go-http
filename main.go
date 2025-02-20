package main

import (
	index_handlers "ducky/http/handlers/index"
	"ducky/http/pkg/server"
)

func main() {
	server := server.NewServer(&server.ServerAddress{Ip: "127.0.0.1", Port: "8008"})
	server.AddHandler("/", "GET", index_handlers.GET)
	server.Start()
}
