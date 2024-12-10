package chat

import (
	"net"
)

type Users map[string]net.Conn
