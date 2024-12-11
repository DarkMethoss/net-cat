package chat

import (
	"net"
	"sync"

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
type BroadcastDetails struct {
	Message string
	User    string
}
type Users map[string]net.Conn
