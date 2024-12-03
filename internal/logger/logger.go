package logger

import (
	"log"
	"os"
)

type Loggers struct {
	LogInfo  *log.Logger
	LogError *log.Logger
	LogChat  *log.Logger
}

func SetLoggers() *Loggers {
	buildLlogFiles()
	err := os.MkdirAll("./logs")
	CheckError(err)
	serverLogFile,err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	CheckError(err)
	chatLogFile,err := os.OpenFile("chat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	logger := new(Loggers) 
	logger.LogInfo,err  := log.New(serverLogFile, os.O_CREATE, 0o666)
	// logger.LogError :=
	// logger.LogChat :=
	return nil
}


buildLogFiles() {

}


func CheckError(err error){
	if err != nil {
		log.Fatal(err)
	}
}