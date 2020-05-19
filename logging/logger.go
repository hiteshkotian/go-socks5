package logging

import "fmt"

type LogWriter interface {
	Write(msg string)
	Close()
}

var logger LogWriter
var loggerChan chan string

func Register(log LogWriter) {
	logger = log
}

func StartLogger() {
	go func() {
		for msg := range loggerChan {
			logger.Write(msg)
		}
	}()
}

func init() {
	loggerChan = make(chan string, 50)
	logger = DefaultLogger{}
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
	var msg string
	if len(parameters) > 0 {
		msg = fmt.Sprintf(message, parameters...)
	}
	msg = fmt.Sprintf("[DEBUG] %s\n", msg)
	loggerChan <- msg
}

func Error(message string, err error, parameters ...interface{}) {
	var msg string
	if len(parameters) > 0 {
		msg = fmt.Sprintf(message, parameters...)
	}
	msg = fmt.Sprintf("[ERROR] %s: %s\n", msg, err.Error())
	loggerChan <- msg
}
