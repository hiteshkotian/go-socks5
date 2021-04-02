package socks5

import "net"

// RequestState type is used to indicate
// the current state a request is in
type RequestState uint8

// AddressType indicates the type of address
// sent in the request
type AddressType uint8

const (
	// RequestStateInit : Initialization State
	RequestStateInit RequestState = 0
	// RequestStateConnecting : Connecting state
	RequestStateConnecting RequestState = 1
	// RequestStateProxying : Proxying request state
	RequestStateProxying RequestState = 2
	// RequestStateTerminating : Terminating connection state
	RequestStateTerminating RequestState = 3
)

// Request holds the properties of a single request
type Request struct {
	State              RequestState // state of the connection
	SourceAddr         net.Addr     // address of the client
	DestinationFQDN    string       // Domain address of the destination. For Socks connect request 0x03
	DestinationAddr    net.Addr     // address of the destination server
	ClientConnection   net.Conn     // Client Connection
	OutboundConnection net.Conn     // Outbound connection
}

// NewRequest creates a new instance of request
func NewRequest(clientConnection net.Conn) *Request {
	request := &Request{}
	request.State = RequestStateInit
	request.ClientConnection = clientConnection
	request.SourceAddr = clientConnection.RemoteAddr()

	return request
}

func (request *Request) Close() error {
	err := request.ClientConnection.Close()
	return err
}
