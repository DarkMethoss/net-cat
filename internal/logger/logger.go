package logger

import (
	"log"
	"os"
)

type Loggers struct {
	LogInfo  *log.Logger
	LogError *log.Logger
}

func SetLoggers() *Loggers {
	serverLogFile := buildLogFiles()
	log.SetFlags(log.Ldate | log.Ltime)
	loggers := new(Loggers)
	loggers.LogInfo = log.New(serverLogFile, "Info: ", log.Ldate|log.Ltime)
	loggers.LogError = log.New(serverLogFile, "Error: ", log.Ldate|log.Ltime)
	return loggers
}

func buildLogFiles() (serverLogFile *os.File) {
	dir := "logs/"
	err := os.MkdirAll(dir, 0o777)
	CheckError(err)
	serverLogFile, err = os.OpenFile(dir+"server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	CheckError(err)
	return serverLogFile
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
