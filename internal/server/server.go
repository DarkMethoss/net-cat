package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"netcat/internal/logger"
)

// var wg
type Server struct {
	Listener net.Listener
	Clients  map[net.Conn]string
	Message  chan string
	Stopped  bool
	mu       sync.Mutex
	*logger.Loggers
}

func NewServer() *Server {
	return &Server{
		Clients: map[net.Conn]string{},
		Message: make(chan string, 1),
		Stopped: false,
		mu:      sync.Mutex{},
		Loggers: logger.SetLoggers(),
	}
}

func (s *Server) Start(port string) {
	fmt.Println("Starting server on port", port)
	var err error
	s.Listener, err = net.Listen("tcp", "localhost:"+port)
	if err != nil {
		s.LogError.Println(err)
		fmt.Println(err)
		return
	}

	for {
		if s.Stopped {
			close(s.Message)
			return // Exit the loop if the server is stopped
		}
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.Stopped {
				s.Message <- "Server is down"
				break // Exit if the server is stopping
			}
			fmt.Println(err)
			continue
		}
		go s.brodcast(conn)
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer s.RemoveClient(conn)
	fmt.Println("New connection from", conn.RemoteAddr())

	conn.Write([]byte("[ENTER YOUR NAME]: "))
	name, _ := bufio.NewReader(conn).ReadString('\n')
	s.mu.Lock()
	s.Clients[conn] = strings.TrimSpace(name)
	s.mu.Unlock()
	s.Message <- fmt.Sprintf("%s has joined the chat\n", s.Clients[conn])
	fmt.Println(s.Clients[conn] + " joined the chat !!!")
	for {
		if s.Stopped {
			s.Message <- fmt.Sprintf("%s has left the chat\n", s.Clients[conn])
			return
		}
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(s.Clients[conn] + " left the chat !!!")
			s.Message <- fmt.Sprintf("%s has left the chat\n", s.Clients[conn])
			s.RemoveClient(conn)
			return
		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		sender := s.Clients[conn]
		prefix := fmt.Sprintf("[%s][%s]:", timestamp, sender)
		s.Message <- fmt.Sprintf("%s%s\n", prefix, message)
	}
}

func (s *Server) RemoveClient(conn net.Conn) {
	s.mu.Lock()
	delete(s.Clients, conn)
	s.mu.Unlock()
}

func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Stopped {
		return
	}
	s.Stopped = true

	// Notify all clients that the server is shutting down
	shutdownMessage := "Server is shutting down...\n"
	for client := range s.Clients {
		client.Write([]byte(shutdownMessage)) // Notify the client
		client.Close()                        // Close the connection
	}
	s.Listener.Close() // Stop accepting new connections
	fmt.Println("Server has been stopped.")
}

func (s *Server) brodcast(conn net.Conn) {
	for message := range s.Message {
		s.mu.Lock()
		for client := range s.Clients {
			if client != conn {
				client.Write([]byte(message))
			}
		}
		s.mu.Unlock()
	}
	fmt.Println("Broadcast loop terminated.")
}
