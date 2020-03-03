package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"hiteshkotian/secure-tcp/server"
)

var serv *server.Server

func main() {
	setupCloseHandler()

	serv, err := server.NewServer(8080)
	if err != nil {
		fmt.Println("Unable to create new server connection: ", err.Error())
		os.Exit(-1)
	}

	err = serv.Start()
	if err != nil {
		fmt.Println("Unable to start tcp server", err.Error())
		os.Exit(1)
	}

	serv.ServeTCP()
	fmt.Println("Returning successfully")
}

// setupCloseHandler will setup a handler to handle
// user interrupt signal to gracefully close all the connections
// and stop the server
func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("shutting down server")
		if serv != nil {
			serv.Stop()
		}
		os.Exit(0)
	}()
}
