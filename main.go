package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	server, err := s.NewTcpChatServer(port)
	if err != nil {
		server.LogError.Println(err)
	} else {
		server.LogInfo.Println("")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	defer func() {
		server.Stopped = true
		server.LogInfo.Println("Shutting down server...")
		server.Listener.Close()
		server.LogInfo.Println("Server stopped gracefully.")
	}()

	// Creation of a socket
	go func() {
		for {
			if server.Stopped {
				return
			}
			conn, err := server.Listener.Accept()
			if err != nil {
				// server.LogError.Println(err)
				continue
			}
			go server.HandleConnection(conn)
		}
	}()

	fmt.Println("Main function starts. Waiting for signal...")
	<-sigChan // Block here until signal is received
	fmt.Println("Signal received. Exiting.")
}
