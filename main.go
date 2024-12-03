package main

import (
	"log"
	"net"

	s "netcat/internal/server"
)

// type chatServer

func main() {
	server := s.NewTcpChatServer()
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		server.LogError.Println("err")
		log.Fatal("")
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			server.LogError.Println(err)
			continue
		}
		go s.HandleConnection(conn)
	}
}
