package handler

import (
	"bufio"
	"net"
	"time"
)

// ProxyState enumeration represents
// the current state of the connection.
type ProxyState int8

const (
	// NEW state is set on initial connection
	NEW ProxyState = iota
	// INITIALIZING state is set while processing
	// the initialize message
	INITIALIZING
	// CONNECTING state is set while trying
	// to establish a connection to the outbound server
	CONNECTING
	// PROXYING state is set when the request is being
	// tunneled from the client to the server
	PROXYING
	// TERMINATING state is set when the inbound and
	// outbound connections are being closed
	TERMINATING
	// ERROR state is set when an error is encountered
	ERROR
)

// Request structure is an abstraction of
// a client request
type Request struct {
	connection net.Conn
	requestID  string

	outboundAddress    string
	outboundIP         string
	outboundPort       uint16
	outboundConnection net.Conn

	state ProxyState

	streamReader *bufio.Reader
	streamWriter *bufio.Writer
}

// NewRequest creates a new instance of request
func NewRequest(conn net.Conn) *Request {
	id := time.Now().String()
	request := Request{requestID: id, connection: conn, state: NEW}

	request.streamReader = bufio.NewReader(conn)
	request.streamWriter = bufio.NewWriter(conn)

	return &request
}

func (request *Request) Read(buffer []byte) (int, error) {
	n, e := request.connection.Read(buffer)
	return n, e
}

func (request *Request) Write(buffer []byte) (int, error) {
	n, e := request.connection.Write(buffer)
	return n, e
}

func (request *Request) SetOutboundIP(ipAddress string) {
	request.outboundIP = ipAddress
}

func (request *Request) SetOutboundPort(port uint16) {
	request.outboundPort = port
}

func (request *Request) SetOutboundConnection(outboundConnetion net.Conn) {
	request.outboundConnection = outboundConnetion
}

func (request *Request) Close() error {
	err := request.connection.Close()
	return err
}

// SetState sets the state of the request
func (request *Request) SetState(state ProxyState) {
	request.state = state
}

// State returns the current state of the request
func (request *Request) State() ProxyState {
	return request.state
}

// Handler interface defines the basic function
// to be implemented by a handler
type Handler interface {
	HandleRequest(*Request) error
}

// HandlerContext object is passed around between
// various handlers to handle the state of the
// request being processed
type HandlerContext struct {
	id int
}

type Proxy struct {
	// Chain of handlers to run
	HandlerChain []Handler
}

func (p *Proxy) ExecuteHandlers(request *Request) {
	for _, handler := range p.HandlerChain {
		handler.HandleRequest(request)
	}
}
