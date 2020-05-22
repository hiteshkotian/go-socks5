package logging

import (
	"fmt"
	"log"
	"os"
)

const channelBufferSize = 50

var logger *log.Logger
var loggerChan chan string
var debugLog bool

func init() {
	logger = log.New(os.Stdout, "Logger :: ", log.Ldate)
	loggerChan = make(chan string, channelBufferSize)
	go func() {
		for msg := range loggerChan {
			logger.Printf(msg)
		}
	}()
}

func EnableDebug() {
	debugLog = true
}

func DisableDebug() {
	debugLog = false
}

func Info(message string, parameters ...interface{}) {
	var msg string
	if len(parameters) > 0 {
		msg = fmt.Sprintf(message, parameters...)
	}
	msg = fmt.Sprintf("[INFO] %s\n", msg)
	loggerChan <- msg
}

func Debug(message string, parameters ...interface{}) {
	if debugLog {
		var msg string
		if len(parameters) > 0 {
			msg = fmt.Sprintf(message, parameters...)
		}
		msg = fmt.Sprintf("[DEBUG] %s\n", msg)
		loggerChan <- msg
	}
}

func Error(message string, err error, parameters ...interface{}) {
	var msg string
	if len(parameters) > 0 {
		msg = fmt.Sprintf(message, parameters...)
	}
	msg = fmt.Sprintf("[ERROR] %s: %s\n", msg, err.Error())
	loggerChan <- msg
}
