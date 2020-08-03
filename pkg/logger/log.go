package logger

import (
	"log"
	"os"
)

var INFO *log.Logger
var ERROR *log.Logger

func InitializeLogger() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	INFO = log.New(file, "INFO: ", log.LstdFlags|log.Lshortfile)
	ERROR = log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile)
}
