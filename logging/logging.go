package logging

import (
	"fmt"
	"log"
	"os"
)

type ProxyLog struct {
	logger *log.Logger
}

var proxyLogger *ProxyLog

func init() {
	proxyLogger = New()
}

func New() *ProxyLog {
	// logger = &Logger{}
	logger := log.New(os.Stdout, "Proxy Logger :: ", log.Ldate)
	return &ProxyLog{logger: logger}
}

func Info(message string, parameters ...interface{}) {
	if len(parameters) > 0 {
		message = fmt.Sprintf(message, parameters)
	}
	proxyLogger.logger.Printf("[INFO] %s\n", message)
}

func Debug(message string, parameters ...interface{}) {

	if len(parameters) > 0 {
		message = fmt.Sprintf(message, parameters)
	}

	proxyLogger.logger.Printf("[DEBUG] %s\n", message)
}

func Error(message string, err error, parameters ...interface{}) {
	if len(parameters) > 0 {
		message = fmt.Sprintf(message, parameters)
	}
	proxyLogger.logger.Printf("[ERROR] %s : %s\n", message, err.Error())
}
