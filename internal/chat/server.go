package chat

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"netcat/internal/helpers"
	"netcat/internal/logger"
)

type Server struct {
	Listener        net.Listener
	Users           Users
	Broadcast       chan BroadcastDetails
	HistoryMessages []string
	Shutdown        bool
	mu              sync.Mutex
	*logger.Loggers
}

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
	s.Listener, err = net.Listen("tcp", "localhost:"+port)
	if err != nil {
		s.LogError.Println("Listening Error", err.Error())
		return
	}
	s.LogInfo.Println("Chat Server Started : server listening for connections on the port " + port)

	go s.brodcast()
	for {
		if s.Shutdown {
			s.Listener.Close() // Stop accepting new connections
			return
		}
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("here", err)
			continue
		}
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	s.LogInfo.Println("New connection from" + conn.RemoteAddr().String())
	userName, err := s.AddUser(conn)
	if err != nil {
		s.LogError.Println("Error", err.Error())
	} else {
		s.LogInfo.Println(userName + " has joined our chat...")
		s.Broadcast <- BroadcastDetails{
			Notification: "\033[32m" + userName + " has joined our chat...\033[0m\n",
			User:         userName,
		}
	}

	defer s.Removeuser(userName)

	err = s.HandleMessages(conn, userName)
	if err != nil {
		s.LogInfo.Println("\033[31m" + userName + " has left our chat...\033[0m\n")
	}
}

func (s *Server) HandleMessages(conn net.Conn, userName string) error {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error Reading message: ", err)
			conn.Close()
			return err
		}

		if strings.TrimSpace(message) != "" {
			s.Broadcast <- BroadcastDetails{
				Message:      helpers.SetPrefix(userName) + message,
				Notification: "",
				User:         userName,
			}
		} else {
			conn.Write([]byte("\033[F\033[2k\r"))
			conn.Write([]byte(helpers.SetPrefix(userName)))
		}

	}
}

func (s *Server) Removeuser(user string) {
	s.mu.Lock()
	delete(s.Users, user)
	defer s.mu.Unlock()
	s.Broadcast <- BroadcastDetails{
		Notification: "\033[31m" + user + " has left our chat...\033[0m\n",
	}
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
				fmt.Println("history:", s.HistoryMessages)
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

func (s *Server) brodcast() {
	for brodcast := range s.Broadcast {
		s.mu.Lock()
		for user, conn := range s.Users {
			if user != brodcast.User && brodcast.Notification == "" {
				conn.Write([]byte("\033[s\n\033[F\033[2K\r"))
				conn.Write([]byte(brodcast.Message))
			} else if user != brodcast.User {
				conn.Write([]byte("\033[s\033[F\033[2K\r"))
				conn.Write([]byte(brodcast.Notification))
			}
			conn.Write([]byte(helpers.SetPrefix(user)))
			if user != brodcast.User {
				conn.Write([]byte("\033[u\033[2B"))
			}
		}
		s.HistoryMessages = append(s.HistoryMessages, brodcast.Message+brodcast.Notification)
		s.mu.Unlock()
	}
}

func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Shutdown = true
	for _, user := range s.Users {
		user.Write([]byte("\nServer has been stopped.")) // Notify the user
		user.Close()                                     // Close the connection
	}
	s.LogInfo.Println("Server has been stoped")
	fmt.Println("Server has been stopped.")
	// Server.Shutdown<-true
}
