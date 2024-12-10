package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"netcat/internal/chat"
	"netcat/internal/helpers"
)

// var logo = ""

func main() {
	port := helpers.HandleArguments()
	if port == "" {
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
	server := chat.NewServer()

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, syscall.SIGINT)
	go server.Start(port)
	defer server.Stop()
	<-stopped
}
