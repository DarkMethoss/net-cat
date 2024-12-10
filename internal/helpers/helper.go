package helpers

import (
	"fmt"
	"os"
	"time"
)

func HandleArguments() (port string) {
	switch len(os.Args[1:]) {
	case 0:
		port = "8989"
	case 1:
		port = os.Args[1]
	default:
		port = ""
	}
	return
}

func SetPrefix(name string) string {
	timestamp := time.Now().Format(time.DateTime)
	return fmt.Sprintf("[%s][%s]:", timestamp, name)
}
