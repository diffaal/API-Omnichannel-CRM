package logger

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *log.Logger

func InitLogger() {
	logFileName := "/var/log/crmx-api.log"
	logFile := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	Logger = log.New(multiWriter, "", log.LstdFlags)
}

func Info(message string) {
	Logger.Println(message)
}

func Error(message string) {
	Logger.Fatalln(message)
}
