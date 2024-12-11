package chat

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"netcat/internal/helpers"
	"netcat/internal/logger"
)

func NewServer() *Server {
	return &Server{
		Users:           Users{},
		Shutdown:        false,
		Broadcast:       make(chan BroadcastDetails),
		HistoryMessages: []string{},
		mu:              sync.Mutex{},
		Loggers:         logger.SetLoggers(),
	}
}

func (s *Server) Start(port string) {
	var err error
	s.Listener, err = net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		s.LogError.Println("Listening Error", err.Error())
		log.Fatal(err)
	}
	s.LogInfo.Println("Chat Server Started : server listening for connections on the port " + port)
	fmt.Println("Chat Server Started : server listening for connections on " + s.Listener.Addr().String())
	go s.brodcast()
	for !s.Shutdown {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("here", err)
			continue
		}
		go s.HandleConnection(conn)
	}
	s.Listener.Close()
}

func (s *Server) HandleConnection(conn net.Conn) {
	if s.Shutdown {
		return
	}
	s.LogInfo.Println("New connection from" + conn.RemoteAddr().String())
	userName, err := s.AddUser(conn)
	if err != nil {
		s.LogError.Println("Error adding user:",)
		return
	}

	// Broadcast user join message
	s.LogInfo.Println(userName + " has joined our chat...")
	s.Broadcast <- BroadcastDetails{
		Message: "\033[32m" + userName + " has joined our chat...\033[0m\n",
		User:    userName,
	}

	defer s.Removeuser(userName)

	// Handle the messages for this user
	s.HandleMessages(conn, userName)
}

func (s *Server) AddUser(conn net.Conn) (string, error) {
	welcomeMessage := "Welcome to TCP-Chat!\n" +
		"         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		" |    .       | ' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     -'       --'\n"
	conn.Write([]byte(welcomeMessage + "\n"))
	for {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		name, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return "", err
		}
		name = name[:len(name)-1]
		if len(s.Users) < 10 {
			err := helpers.ValidName(name)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			if _, exists := s.Users[name]; !exists {
				s.mu.Lock()
				s.Users[name] = conn
				s.mu.Unlock()
				conn.Write([]byte(strings.Join(s.HistoryMessages, "")))
				return name, nil
			}
			conn.Write([]byte("UserName Already Exists\n"))
			continue
		} else {
			conn.Write([]byte("Server Full"))
			return "", errors.New("Server Full")
		}
	}
}

func (s *Server) HandleMessages(conn net.Conn, userName string) {
	for {
		if s.Shutdown {
			conn.Close()
			break
		}

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			// Log the error
			s.LogInfo.Println(userName + " has left the chat (EOF)")
			return
		}
		message, err = helpers.ValidMessage(strings.TrimSpace(message))
		if err == nil && message != "" {
			s.Broadcast <- BroadcastDetails{
				Message: helpers.SetPrefix(userName) + message + "\n",
				User:    userName,
			}
		} else if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			conn.Write([]byte(helpers.SetPrefix(userName)))
		} else {
			conn.Write([]byte(helpers.SetPrefix(userName)))
		}
	}
}

func (s *Server) Removeuser(user string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Broadcast <- BroadcastDetails{
		Message: "\033[31m" + user + " has left our chat...\033[0m\n",
	}
	s.Users[user].Close()
	delete(s.Users, user)
}


func (s *Server) brodcast() {
	for brodcast := range s.Broadcast {
		s.mu.Lock()
		for user, conn := range s.Users {
			if user != brodcast.User {
				conn.Write([]byte("\033[s\n\033[F\033[2K\r"))
				conn.Write([]byte(brodcast.Message))
			}
			conn.Write([]byte(helpers.SetPrefix(user)))
			if user != brodcast.User {
				conn.Write([]byte("\033[u\033[B"))
			}
		}
		s.HistoryMessages = append(s.HistoryMessages, brodcast.Message)
		s.mu.Unlock()
	}
}

func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Shutdown = true
	close(s.Broadcast)
	for _, user := range s.Users {
		user.Write([]byte("\nServer has been stopped.")) // Notify the user
	}
	s.LogInfo.Println("Server has been stoped")
	fmt.Println("Server has been stopped.")
}
