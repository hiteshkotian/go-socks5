package proxy

import (
	"context"
	"fmt"
	"hiteshkotian/ssl-tunnel/handler"
	"hiteshkotian/ssl-tunnel/logging"
	"hiteshkotian/ssl-tunnel/socks5"
	"net"
	"time"
)

type ctxKey string

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
	logging.Info("Starting Proxy Server")
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

		if err != nil {
			logging.Error("Erorr while reading incoming request", err)
			return err
		}

		logging.Debug("received connection request from %s",
			conn.RemoteAddr().String())

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
				return
			}
			server.sem <- true

			ctx := context.Background()
			go server.handleRequest2(ctx, conn, server.sem)
		}
	}
}

func (server *Server) handleRequest2(ctx context.Context,
	conn net.Conn, sem chan bool) {
	logging.Debug("Processing incoming client request")

	request := socks5.NewRequest(conn)

	processRequest := true
	for processRequest {
		// Step 1 : Handle Initiial
		switch request.State {
		case socks5.RequestStateInit:
			server.handleInitialLocal(request)
		case socks5.RequestStateConnecting:
			server.handleConnectLocal(request)
		case socks5.RequestStateProxying:
			server.startProxying(request)
		case socks5.RequestStateTerminating:
			request.Close()
			<-sem
			processRequest = false
		}
	}
}

func (server *Server) startProxying(request *socks5.Request) {
	outboundHandler := handler.OutboundHandler{}
	err := outboundHandler.HandleRequest(request)
	if err != nil {
		request.State = socks5.RequestStateTerminating
		return
	}

	request.State = socks5.RequestStateTerminating
}

func (server *Server) handleInitialLocal(request *socks5.Request) {
	// Initial request structre is :
	// init_request_pkt {
	// 		version (1) = 0x05
	//		nmethods (1)
	//		methods (1...255)
	// }
	// Response structure is :
	// init_response_pkt {
	// 		version (1) = 0x05
	//		method (1)
	// }
	// TODO Check how to handle authentication request
	logging.Debug("Processing init request")
	clientConn := request.ClientConnection

	requestStream := make([]byte, 260)

	n, e := clientConn.Read(requestStream)
	if e != nil {
		logging.Error("Error reading response", e)
	} else if n < 2 {
		logging.Error("Insufficient bytes read", nil)
	}

	logging.DumpHex(requestStream[:n], "INIT Method")

	_, err := socks5.GetSocketInitialSerialized(requestStream[:n])

	if err != nil {
		response, _ := socks5.GetSocketInitialResponseSerialized(0xFF)
		clientConn.Write(response)
		request.State = socks5.RequestStateTerminating
		return
	}

	// TODO Add support for authentication

	response, _ := socks5.GetSocketInitialResponseSerialized(0x00)
	logging.DumpHex(response, "Sending response")
	clientConn.Write(response)
	// Change the state
	request.State = socks5.RequestStateConnecting
}
func (server *Server) handleConnectLocal(request *socks5.Request) {
	// Connect request format
	// connect_req_pkt {
	//		version (1) = 0x05
	//		command (1) = [0x01, 0x02, 0x03]
	//		rsv (1) = 0x00
	//		atyp (1) = [0x01, 0x03, 0x04]
	//		dst.addr (...)
	//		dst.port (2)
	// }
	// Connect response format
	// connect_response_pkt {
	//		version (1) = 0x05
	//		rep (1) = [0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09]
	//		rsv (1) = 0x00
	//		atyp (1) = [0x01, 0x03, 0x04]
	//		bind.addr (...)
	//		bind.port (2)
	// }

	clientConn := request.ClientConnection
	requestStream := make([]byte, 512)

	n, e := clientConn.Read(requestStream)
	if e != nil || n <= 0 {
		// Send error
		response, _ := socks5.GetSocketInitialResponseSerialized(0xFF)
		clientConn.Write(response)
		request.State = socks5.RequestStateTerminating
		return
	}

	connectRequest, err := socks5.GetSocketRequestDeserialized(requestStream[:n])

	if err != nil {
		response, _ := socks5.GetSocketInitialResponseSerialized(0xFF)
		clientConn.Write(response)
		request.State = socks5.RequestStateTerminating
		return
	}

	reply := socks5.CreateSocksReply(connectRequest)

	// Create connection
	request.OutboundConnection, err = server.createOuboundConnection(connectRequest)
	if err != nil {
		logging.Error("Error connecting to remote host", err)
		reply.SetReply(socks5.ReplyNetUnreachable)
		request.State = socks5.RequestStateTerminating
		// return
	} else {
		request.State = socks5.RequestStateProxying
	}

	replyStream, _ := socks5.GetSocketResponseSerialized(reply)
	clientConn.Write(replyStream)
	return
}

func (server *Server) createOuboundConnection(
	connectRequest socks5.SockRequest) (outConnection net.Conn, err error) {

	var address string
	var network string

	switch connectRequest.GetAddressType() {
	case socks5.AtypIPV4:
		network = "tcp"
		ip := net.IP(connectRequest.GetDestinationAddress())
		address = fmt.Sprintf("%s:%d", ip.String(), connectRequest.GetDestinationPort())
	case socks5.AtypIPV6:
		network = "tcp6"
		ip := net.IP(connectRequest.GetDestinationAddress())
		address = fmt.Sprintf("[%s]:%d", ip.String(), connectRequest.GetDestinationPort())
		// TODO Lookup Domain
	}

	outConnection, err = net.Dial(network, address)
	return
}

// Stop stops the server
func (server *Server) Stop() {
	// Closing Channel
	logging.Info("Stopping Proxy Server")
	close(server.connectHandler)
	server.listener.Close()
}
