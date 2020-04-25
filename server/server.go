package server

import (
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

const defaultReadTTL = 10
const defaultWriteTTL = 30

//Server structure representing server properties
type Server struct {
	port     int
	readTTL  time.Duration
	writeTTL time.Duration
	listener net.Listener
}

// NewServer creates a new instance of the server
func NewServer(port int) (*Server, error) {
	if port <= 0 {
		return nil, fmt.Errorf("Invalid port number : %d", port)
	}
	s := Server{
		port:     port,
		readTTL:  defaultReadTTL * time.Second,
		writeTTL: defaultWriteTTL * time.Second,
	}
	return &s, nil
}

// Start starts the server and accepts incoming client requests
func (server *Server) Start() error {
	var err error
	fmt.Println("Starting Server")
	server.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", server.port))
	if err != nil {
		fmt.Println("Unable to start tcp server: ", err.Error())
		return err
	}
	return server.ServeTCP()
}

// ServeTCP will start  accepting TCP connections and will
// respond to client's connection.
func (server *Server) ServeTCP() error {
	for {
		conn, err := server.listener.Accept()
		fmt.Printf("received connection request from %s\n",
			conn.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Set all the required timeouts
		conn.SetReadDeadline(
			time.Now().Add(server.readTTL))
		conn.SetWriteDeadline(
			time.Now().Add(server.writeTTL))

		// Create a new thread to handle incoming connections
		go handleConnection(conn)
	}
}

// Stop will stop the tcp server
// and close all existing client connections
func (server *Server) Stop() {
	fmt.Println("stopping server")
	server.listener.Close()
}

func isHTTPMessage(msg string) bool {
	reg := "^(GET|PUT|POST|DELETE|CONNECT).*"
	val, err := regexp.MatchString(reg, msg)
	if err != nil {
		fmt.Println("Error when running regexp: ", err.Error())
	}
	fmt.Printf("Val is : %t\n", val)
	return val
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	msgData := make([]byte, 2048)
	var len int
	var err error
	len, err = conn.Read(msgData)

	if err != nil {
		fmt.Println("error encountered while reading", err)
		conn.Close()
		return
	}

	fmt.Printf("Data read from the client : %d\n", len)

	msg := string(msgData[:len])

	if isHTTPMessage(msg) {
		fmt.Println("We have HTTP Message in our hand")
		fmt.Println("Setting up a tunnel")

		msgSplit := strings.Split(string(msgData), "\r\n\r\n")
		fmt.Println("\n\n\nData is : ", msgSplit[1])
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nHost: localhost\r\n\r\nHello from the server.\nI got your message\r\n")))
		return
	}

	trimmedMsg := strings.TrimSpace(msg)
	trimmedMsg = strings.TrimSuffix(trimmedMsg, "\r\n")

	hexString := hex.EncodeToString([]byte(trimmedMsg))
	fmt.Printf("Command from Client: <<<<<<< %s >>>>>>\n", hexString)

	if trimmedMsg == "HELLO" {
		conn.Write([]byte(fmt.Sprintf("Hello from %s\n", conn.LocalAddr().String())))
		return
	} else if trimmedMsg == "TIME" {
		conn.Write([]byte(fmt.Sprintf("Hello from %s\n", time.Now().String())))
	} else {
		conn.Write([]byte(fmt.Sprintf("invalid command :%s\n", trimmedMsg)))
	}
}
