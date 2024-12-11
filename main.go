package main

import (
	"fmt"
	"os"
	"os/signal"

	"netcat/internal/chat"
	"netcat/internal/helpers"
)

// var logo = ""

func main() {
	port := helpers.HandleArguments()
	if port == "" {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	server := chat.NewServer()

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, os.Interrupt)
	go server.Start(port)
	defer server.Stop()
	<-stopped
}
