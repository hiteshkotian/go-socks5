package proxy

import (
	"fmt"
	"hiteshkotian/ssl-tunnel/handler"
	"hiteshkotian/ssl-tunnel/logging"
	"net"
	"time"
)

// Server structure represents the main proxy instance
type Server struct {
	// Name of the server
	name string
	// Port the server is listening for incoming requests
	port int
	// Maximum number of concurrent connections
	// that can be processed at a given time
	maxConnectionCount int
	// incoming network listener
	listener net.Listener
	// connectionHandler channel. This channel is used for piping the
	// incoming connections to the appropriate handler
	connectHandler chan net.Conn
	// Connection limiter. This channel ensures that at a given time the
	// configured number of requests are being processed.
	sem chan bool
}

// Request represents an incoming request from a client
type Request struct {
	// Unique ID assigned to the request
	requestID int
	// Client network connection object
	connection net.Conn
	// Client address
	remoteAddress string
}

// New creats a new instance of the proxy
func New(name string, port, maxConnectionCount int) *Server {

	proxy := &Server{name: name, port: port,
		maxConnectionCount: maxConnectionCount}
	proxy.connectHandler = make(chan net.Conn)
	proxy.sem = make(chan bool, proxy.maxConnectionCount)

	return proxy
}

// NewFromConfig reads the provided config file and
// returns a proxy instance
func NewFromConfig(configPath string) (*Server, error) {
	// TODO to implement this
	return nil, nil
}

// Start starts the server and accepts incoming client requests
func (server *Server) Start() error {
	var err error
	logging.Debug("Starting Proxy Server")
	server.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", server.port))
	if err != nil {
		logging.Error("Unable to start tcp server: ", err)
		return err
	}

	return server.ServeTCP()
}

// ServeTCP will start  accepting TCP connections and will
// respond to client's connection.
func (server *Server) ServeTCP() error {

	// Start the connection handler
	go server.startHandler()

	for {
		conn, err := server.listener.Accept()
		logging.Info("received connection request from %s",
			conn.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Set all the required timeouts
		conn.SetReadDeadline(
			time.Now().Add(10 * time.Second))
		conn.SetWriteDeadline(
			time.Now().Add(30 * time.Second))

		server.connectHandler <- conn
	}
}

// startHandler function starts listening for incoming TCP
// connection and handles the incoming requests
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
			go server.handleRequest(conn, server.sem)
		default:
		}
	}
}

// HandleRequest implementation for TCP Handler.
// This function will accept all incoming TCP requests
// and serialize it to a request object if it is a valid request.
// In case of a serialization issue, the handler will return
// an appropriate error code to the client.
func (server *Server) handleRequest(conn net.Conn, sem chan bool) error {
	fmt.Println("Processing incoming request in tcp handler")
	defer conn.Close()

	request := handler.NewRequest(conn, []byte("msgData"))
	handlerChain := [2]handler.Handler{&handler.CustomHandler{}, &handler.OutboundHandler{}}

	for _, handler := range handlerChain {
		err := handler.HandleRequest(request)
		if err != nil {
			fmt.Printf("Error while sending request : %s\n", err.Error())
			<-sem
			return err
		}
		// }
		<-sem
	}
	// outboundHandler := handler.OutboundHandler{}
	// err := outboundHandler.HandleRequest(request)
	// if err != nil {
	// 	fmt.Printf("Error while sending request : %s\n", err.Error())
	// 	<-sem
	// 	return err
	// }
	// // }
	// <-sem
	return nil
}

// Stop stops the server
func (server *Server) Stop() {
	// Closing Channel
	fmt.Println("Stopping Proxy Server")
	close(server.connectHandler)
	server.listener.Close()
}
