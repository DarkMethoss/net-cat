package server

import (
	"fmt"
	"io"
	"net"

	"netcat/internal/logger"
)

type TcpChatServer struct {
	Listener net.Listener
	clients  map[net.Conn]string
	*logger.Loggers
}

func NewTcpChatServer(port string) (*TcpChatServer, error) {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		return nil, fmt.Errorf("failed start server at localhost:%w", err)
	}

	return &TcpChatServer{
		Listener: listener,
		clients:  map[net.Conn]string{},
		Loggers:  logger.SetLoggers(),
	}, nil
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
