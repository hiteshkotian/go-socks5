package main

import (
	"fmt"
	"hiteshkotian/ssl-tunnel/proxy"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	proxy := proxy.New("server1", 8080)
	setupCloseHandler(proxy)

	proxy.Start()
}

func setupCloseHandler(proxy *proxy.Server) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("shutting down proxy server")
		once := sync.Once{}
		onceBody := func() {
			proxy.Stop()
		}
		once.Do(onceBody)
		os.Exit(0)
	}()
}
