// Package main contains the main entry point for the proxy
package main

import (
	"fmt"
	"hiteshkotian/ssl-tunnel/proxy"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Main entry point of the proxy
func main() {
	// Set the proxy properties
	// TODO either accept flags or configuration file
	// to bootstrap the proxy
	name := "server1"
	port := 8080
	maxConnCount := 2
	// Create an instance of the proxy
	proxy := proxy.New(name, port, maxConnCount)

	// Setup the close handlers to handle interrupts
	setupCloseHandler(proxy)

	// Start the proxy
	proxy.Start()
}

// setupCloseHandler function registers SIGTERM signal
// to gracefully shutdown the server
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
