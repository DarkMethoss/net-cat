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
}

func StartTcpChatServer() *TcpChatServer {
	server := &TcpChatServer{
		users:   map[net.Conn]string{},
		Loggers: logger.SetLoggers(),
	}
	return server
}

func (tcs *TcpChatServer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from", conn.RemoteAddr())

	// Echo the received data back to the client

	// for conn, name := range tcs.users {

	if _, err := io.Copy(conn, conn); err != nil {
		tcs.LogError.Println(err)
	}
	// }
}


