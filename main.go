package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"netcat/internal/helpers"
	s "netcat/internal/server"
)

// ! strconv not allowed
// ! check port range & port 0 not allowed

// var logo = ""

func main() {
	port := helpers.HandleArguments()
	if port == "" {
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
	server := s.NewServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	// defer func() {
	// 	server.Stopped = true
	// 	server.LogInfo.Println("Shutting down server...")
	// 	server.LogInfo.Println("Server stopped gracefully.")
	// }()
	defer server.Stop()
	go server.Start(port)
	<-sigChan
}
