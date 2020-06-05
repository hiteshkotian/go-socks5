// Package main contains the main entry point for the proxy
package main

import (
	"hiteshkotian/ssl-tunnel/logging"
	"hiteshkotian/ssl-tunnel/proxy"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var version string

// Main entry point of the proxy
func main() {
	// logger := logging.New()
	logging.Info("Iniitializing proxy tunnel version: %s", version)
	// fmt.Printf("Project version is %s\n", version)

	// Set the proxy properties
	// TODO either accept flags or configuration file
	// to bootstrap the proxy
	name := "server1"
	port := 1080
	maxConnCount := 3
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
		logging.Info("shutting down proxy server")
		once := sync.Once{}
		onceBody := func() {
			proxy.Stop()
		}
		once.Do(onceBody)
		os.Exit(0)
	}()
}
