package main

import (
	"fmt"
	"hiteshkotian/secure-tcp/proxy"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	setupCloseHandler()

	proxy, err := proxy.DaemonNew("id", "")
	if err != nil {
		fmt.Println("Error when creating proxy daemon", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Proxy Daemon : \nID: %s\nProxyID : %s\n", proxy.GetID(), proxy.GetProxyID())
}

func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("shutting down server")
		os.Exit(0)
	}()
}
