package main

import (
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8989"
	SERVER_TYPE = "tcp"
)

func main() {
	listner, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

}


func checkError(err error) {
	
}