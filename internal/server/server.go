package server

import (
	"fmt"
	"io"
	"net"

	"netcat/internal/logger"
)

type TcpChatServer struct {
	users map[net.Conn]string
	*logger.Loggers
	chatLog.
}

func NewTcpChatServer() *TcpChatServer {
	server := &TcpChatServer{
		users:   map[net.Conn]string{},
		Loggers: logger.SetLoggers(),
	}
	server.LogInfo.Println("Chat server started.")
	return server
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from", conn.RemoteAddr())

	// Echo the received data back to the client
	if _, err := io.Copy(conn, conn); err != nil {
		fmt.Println("Error echoing data:", err)
	}
}
