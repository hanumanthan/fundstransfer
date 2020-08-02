package logger

import (
	"log"
	"os"
)

var INFOLOG *log.Logger
var ERRORLOG *log.Logger

func InitializeLogger() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	INFOLOG = log.New(file, "INFO: ", log.LstdFlags|log.Lshortfile)
	ERRORLOG = log.New(file, "ERROR: ", log.Ldate|log.Lshortfile)
}
