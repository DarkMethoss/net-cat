package main

import (
	"fmt"
	"log"
	"net"
	"os"

	s "netcat/internal/server"
)

// ! strconv not allowed
// ! check port range & port 0 not allowed

func main() {
	port := ""

	switch len(os.Args[1:]) {
	case 0:
		port = "8989"
	case 1:
		port = os.Args[1]
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
	}

	server := s.StartTcpChatServer()
	listener, err := net.Listen("tcp", "localhost:"+port)

	if err != nil {
		server.LogError.Println("err")
		log.Fatal("")
	} else {
		server.LogInfo.Println("Server Listening to port " + port)
	}

	stop := make(chan struct{})
	defer func() {
		fmt.Println("Server has been stoped")
		server.LogInfo.Println("Server Stoped")
		listener.Close()
	}()

	// Creation of a socket
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				server.LogError.Println(err)
				continue
			}
			go server.HandleConnection(conn)
		}
	}()
	<-stop
}
