package handler

import (
	"hiteshkotian/ssl-tunnel/socks5"
)

// Handler interface defines the basic function
// to be implemented by a handler
type Handler interface {
	HandleRequest(*socks5.Request) error
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

func (p *Proxy) ExecuteHandlers(request *socks5.Request) {
	for _, handler := range p.HandlerChain {
		handler.HandleRequest(request)
	}
}
