package main

import (
	"ducky/http/handlers/index"
	"ducky/http/pkg/server"
)

func main() {
	server := server.NewServer(&server.ServerAddress{Ip: "127.0.0.1", Port: "8008"})
	server.AddHandler("/", "GET", index.GET)
	server.Start()
}
