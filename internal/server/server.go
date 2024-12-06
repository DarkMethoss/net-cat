package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"

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
		fmt.Println(err)
		return
	}

	for {
		if s.Stopped {
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
		go s.HandleConnection(conn)
		go s.brodcast(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer func() {
		conn.Write([]byte("Server Down"))
		s.mu.Lock()
		delete(s.Clients, conn)
		s.mu.Unlock()
		conn.Close()
	}()
	fmt.Println("New connection from", conn.RemoteAddr())

	conn.Write([]byte("[ENTER YOUR NAME]: "))
	name, _ := bufio.NewReader(conn).ReadString('\n')

	s.Clients[conn] = name
	fmt.Println(s.Clients)

	for {
		if s.Stopped {
			s.Message <- fmt.Sprintln("Server Down !!")
			s.mu.Lock()
			delete(s.Clients, conn)
			s.mu.Unlock()
			return
		}
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			s.mu.Lock()
			delete(s.Clients, conn)
			s.mu.Unlock()
			s.Message <- fmt.Sprintf("%s has left the chat\n", s.Clients[conn])
			return
		}

		s.Message <- fmt.Sprintf("%s: %s", s.Clients[conn], message)
	}
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
	for {
		message := <-s.Message
		for client := range s.Clients {
			if client != conn {
				client.Write([]byte(message))
			}
		}
	}
}
