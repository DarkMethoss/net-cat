package main

import (
	"fmt"
	"os"
	"os/signal"

	s "netcat/internal/server"
)

// ! strconv not allowed
// ! check port range & port 0 not allowed

type Signal string

func (sig Signal) String() string {
	switch sig {
	case "0x2":
		return "SIGINT"
	default:
		return "UNKNOWN_SIGNAL"
	}
}

func (sig Signal) Signal() {
	// Marker method
}

const (
	SIGINT = Signal("0x2") // Custom SIGINT representation
)

func main() {
	port := ""
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Signal(SIGINT))
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
	}

	stop := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	defer func() {
		server.LogInfo.Println("Shutting down server...")
		server.Listener.Close()
		close(stop)
		server.LogInfo.Println("Server stopped gracefully.")
	}()

	// Creation of a socket
	go func() {
		for {
			conn, err := server.Listener.Accept()
			if err != nil {
				server.LogError.Println(err)
				continue
			}
			go server.HandleConnection(conn)
		}
	}()
	
	<-sigChan
}
