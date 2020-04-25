package handler

import "net"

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
	HandleRequest(net.Conn) error
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

func (p *Proxy) ExecuteHandlers() {
	for _, handler := range p.HandlerChain {
		handler.HandleRequest(nil)
	}
}
