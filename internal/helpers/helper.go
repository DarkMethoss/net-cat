package helpers

import (
	"errors"
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

func ValidName(name string) error {
	if len(name) < 3 {
		return errors.New("name should contain at least have 3 characters")
	}
	if len(name) > 13 {
		return errors.New("name should contain less than 13 characters")
	}
	err := errors.New("name should only contain (alphanumerical, _, -) characters ")
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-') {
			return err
		}
	}
	return nil
}

func ValidMessage(message string) (string, error) {
	if message == "" {
		return "", nil
	}
	for _, char := range message {
		if char < 32 || char > 126 {
			return "", errors.New("message should only contain printable ascii characters")
		}
	}
	return message, nil
}
