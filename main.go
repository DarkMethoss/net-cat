package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

// type chatServer

func main() {
	listner, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("listening to port: 8080")
	}

	conn, err := listner.Accept()
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte("Welcome to the chat."))
	if err != nil {
		fmt.Println("Error writing to the connection: ", err)
	}

	proxyConn(listner.Addr().String(), conn.LocalAddr().String())
}

func proxyConn(source, destination string) error {
	connSource, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}
	defer connSource.Close()

	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer connDestination.Close()
	go func() { _, _ = io.Copy(connSource, connDestination) }()
	_, err = io.Copy(connDestination, connSource)
	return err
}
