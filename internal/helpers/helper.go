package helpers

import (
	"os"
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
