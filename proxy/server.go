package proxy

import (
	"fmt"
	"hiteshkotian/ssl-tunnel/handler"
	"net"
	"time"
)

// Server structure represents the main proxy instance
type Server struct {
	name           string
	port           int
	listener       net.Listener
	connectHandler chan net.Conn
	sem            chan bool
}

// New creats a new instance of the proxy
func New(name string, port int) *Server {
	proxy := &Server{name: name, port: port}

	proxy.connectHandler = make(chan net.Conn)
	proxy.sem = make(chan bool, 2)

	return proxy
}

// NewFromConfig reads the provided config file and
// returns a proxy instance
func NewFromConfig(configPath string) (*Server, error) {
	return nil, nil
}

// Start starts the server and accepts incoming client requests
func (server *Server) Start() error {
	var err error
	fmt.Println("Starting Proxy Server")
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

	go server.startHandler()

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
			time.Now().Add(30 * time.Second))
		conn.SetWriteDeadline(
			time.Now().Add(60 * time.Second))
		conn.SetDeadline(time.Now().Add(1 * time.Second))

		server.connectHandler <- conn
	}
}

func (server *Server) startHandler() {
	for {
		select {
		case conn, more := <-server.connectHandler:
			if !more {
				fmt.Println("Closing")
				return
			}
			fmt.Println("--------> Received request")
			server.sem <- true
			handler := handler.TcpHandler{}
			handler.HandleRequest(conn, server.sem)
		default:
		}
	}
}

// Stop stops the server
func (server *Server) Stop() {
	// Closing Channel
	fmt.Println("Stopping Proxy Server")
	close(server.connectHandler)
	server.listener.Close()
}
