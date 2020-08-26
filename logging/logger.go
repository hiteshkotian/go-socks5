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

func DumpHex(stream []byte, message string, parameters ...interface{}) {
	// if !debugLog {
	// 	return
	// }

	// Generate the hex dump
	streamLen := len(stream)
	var hexStream string
	for _, val := range stream {
		if len(hexStream) == 0 {
			hexStream = fmt.Sprintf("0x%02x", val)
		} else {
			hexStream = fmt.Sprintf("%s 0x%02x", hexStream, val)
		}
		// if index%16 == 0 {
		// 	hexStream = fmt.Sprintf("%s\n", hexStream)
		// }
	}

	var msg string
	if len(parameters) > 0 {
		msg = fmt.Sprintf(message, parameters...)
	} else {
		msg = message
	}
	msg = fmt.Sprintf("[DEBUG] [len : %d] %s :: [%s]\n", streamLen, msg, hexStream)
	loggerChan <- msg
}

func Info(message string, parameters ...interface{}) {
	var msg string
	if len(parameters) > 0 {
		msg = fmt.Sprintf(message, parameters...)
	} else {
		msg = message
	}
	msg = fmt.Sprintf("[INFO] %s\n", msg)
	loggerChan <- msg
}

func Debug(message string, parameters ...interface{}) {
	if debugLog {
		var msg string
		if len(parameters) > 0 {
			msg = fmt.Sprintf(message, parameters...)
		} else {
			msg = message
		}
		msg = fmt.Sprintf("[DEBUG] %s\n", msg)
		loggerChan <- msg
	}
}

func Error(message string, err error, parameters ...interface{}) {
	var msg string
	if len(parameters) > 0 {
		msg = fmt.Sprintf(message, parameters...)
	} else {
		msg = message
	}

	if err == nil {
		msg = fmt.Sprintf("[ERROR] %s\n", msg)
	} else {
		msg = fmt.Sprintf("[ERROR] %s: %s\n", msg, err.Error())
	}

	loggerChan <- msg
}
