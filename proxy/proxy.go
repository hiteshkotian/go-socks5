// Package proxy will act as the intermediary between the
// client and the connecting server
// Proxy has the following state :
// Initializing -> Ready -> listening -> validating client -> connecting -> connected -> proxying -> terminating | timedout -> disconnected -> ready (back)
// The proxy can be a TCP proxy, and do DNS lookups for the client
package proxy

import (
	"fmt"
	"os"
)

// State represents the state the proxy is in
type State int64

// States the proxy will be in
const (
	// State when the proxy connection is starting up.
	// At this point the proxy will configure itself when starting up,
	// or will deallocated all the resources it allocated when it was
	// processing the last request
	Initializing State = 0
	// State when the proxy connection is done initializing and
	// is now ready to accept client connection
	Ready State = 1
	// State when the proxy gets a client connection. At this state,
	// the proxy will validate the request, check if the packet is formatted
	// correctly and will also check the authentication token
	// (if the request is a proxy request). Also the proxy will
	// check if the server the client wishes to connec to is reachable.
	Connecting State = 2
	// After successfully validating the request in the Connecting state,
	// the proxy is now ready to proxy connection.
	Proxying State = 3
	// After the connection is passed, the proxy will go in this state
	// at which point it will stop accepting requests and gracefully
	// close the connection and any resources
	Terminating State = 4
	// State the proxy connection goes to when the client connection timed out.
	// This could also happen when the server does not respond in time.
	Timeout State = 5
	// State the proxy is in when it is disconnected from the client
	// and is in the process of clearning up connections.
	// TODO : Maybe we don't need this
	Disconnected State = 6
)

// Daemon represents a single proxy connection
type Daemon struct {
	// Current state of the proxy
	state State
	// ID of the daemon.
	// This will help us uniquely identify which thread
	// Processed a request
	id string
	// ID of the proxy, will be used by the client to know which proxy
	// setup it is connected to.
	// The proxy ID can be configured as a regular expression,
	// by default it will be the hostname.
	proxyID string
}

// DaemonNew will create a new instance of the daemon
func DaemonNew(id, proxyID string) (*Daemon, error) {
	// Check ID
	if id == "" {
		return nil, fmt.Errorf("Invalid ID provided")
	}
	// Check ProxyID
	if proxyID == "" {
		var err error
		proxyID, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}
	return &Daemon{id: id, proxyID: proxyID, state: Initializing}, nil
}

// GetState returns the current state of the daemon
func (d *Daemon) GetState() State {
	return d.state
}

// GetID returns the ID of the daemon
func (d *Daemon) GetID() string {
	return d.id
}

// GetProxyID returns the proxy ID of the daemon
func (d *Daemon) GetProxyID() string {
	return d.proxyID
}
