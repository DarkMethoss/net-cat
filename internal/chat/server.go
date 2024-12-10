package chat

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
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
	Shutdown        chan os.Signal
	mu              sync.Mutex
	*logger.Loggers
}

func NewServer() *Server {
	return &Server{
		Users:           Users{},
		Shutdown:        make(chan os.Signal, 1),
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
		s.Log("Info", err.Error())
		return
	}
	s.Log("Info", "Chat Server Started : server listening for connections on the port "+port)

	go s.brodcast()
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("here", err)
			continue
		}
		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	s.Log("Info", "New connection from"+conn.RemoteAddr().String())
	userName, err := s.AddUser(conn)
	if err != nil {
		s.Log("Error", err.Error())
	} else {
		s.Log("Infor", userName+" has joined our chat...\n")
		s.Broadcast <- BroadcastDetails{
			Notification: "\n" + userName + " has joined our chat...\n",
			User:         userName,
		}
	}

	defer s.Removeuser(userName)

	err = s.HandleMessage(conn, userName)
	if err != nil {
		s.Log("Infor", userName+" has left our chat...")
	}
}

func (s *Server) HandleMessage(conn net.Conn, userName string) error {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error Reading message: ", err)
			s.Broadcast <- BroadcastDetails{
				Message:      "",
				Notification: "\n" + userName + " has left our chat...",
				User:         userName,
			}
			return err
		}

		if strings.TrimSpace(message) != "" {
			s.Broadcast <- BroadcastDetails{
				Message:      message,
				Notification: "",
				User:         userName,
			}
		} else {
			conn.Write([]byte(helpers.SetPrefix(userName)))
		}

	}
}

func (s *Server) Removeuser(user string) {
	s.mu.Lock()
	delete(s.Users, user)
	s.mu.Unlock()
}

func (s *Server) Stop() {
	s.mu.Lock()
	s.Listener.Close() // Stop accepting new connections
	defer s.mu.Unlock()
	for _, user := range s.Users {
		user.Write([]byte("\nServer has been stopped.")) // Notify the user
		user.Close()                                     // Close the connection
	}
	fmt.Println("Server has been stopped.")
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
		fmt.Println([]byte(name))
		if len(s.Users) < 10 {
			if _, exists := s.Users[name]; !exists {
				s.mu.Lock()
				s.Users[name] = conn
				s.mu.Unlock()
				return name, nil
			}
			return "", errors.New("UserName Slready Exists")
		} else {
			conn.Write([]byte("Server Full"))
			return "", errors.New("Server Full")
		}
	}
}

func (s *Server) brodcast() {
	for brodcast := range s.Broadcast {
		s.mu.Lock()
		for user := range s.Users {
			if user != brodcast.User && brodcast.Notification == "" {
				s.Users[user].Write([]byte("\n" + helpers.SetPrefix(brodcast.User)))
				s.Users[user].Write([]byte(brodcast.Message))
			} else if user != brodcast.User && brodcast.Notification != "" {
				s.Users[user].Write([]byte(brodcast.Notification))
			}
			s.Users[user].Write([]byte(helpers.SetPrefix(user)))
		}
		s.mu.Unlock()
	}
}

func (s *Server) Log(logType string, logMessage string) {
	switch logType {
	case "Info":
		s.LogInfo.Println(logMessage)
	case "Error":
		s.LogError.Println(logMessage)
	}
}
