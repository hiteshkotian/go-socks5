package handler

import (
	"net"
	"time"
)

type Request struct {
	id   int64
	conn net.Conn
	data []byte
	host string
}

func NewRequest(conn net.Conn, data []byte) *Request {
	id := time.Now()
	request := Request{id: id.Unix(), conn: conn, host: "www.google.com:443", data: data}
	return &request
}

// Handler interface defines the basic function
// to be implemented by a handler
type Handler interface {
	// NewHandler will return a new handler instance
	// NewHandler() (Handler, error)
	// NewHandlerWithConfig will return a new handler
	// instance based on the config provided
	// NewHandlerWithConfig(map[string]string) (Handler, error)
	// HandleRequest will handle the incoming request
	// and subsequently handle the response
	// HandleRequest(net.Conn) error

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
