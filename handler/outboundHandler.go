package handler

import (
	"hiteshkotian/ssl-tunnel/logging"
	"hiteshkotian/ssl-tunnel/socks5"
	"net"
	"time"
)

type OutboundHandler struct {
}

// proxyData function will read the data from the "from" channel
// and synchronously write it to the "to" channel.
// On operation complete the function will write true to the done and the
// complete channel.
// If the other read go routine is done, the signal will be received in the
// otherDone channel, at which point the function will stop all action.
// NOTE: We could make this to stop on read "\r\n\r\n" but then we are just
// delimiting it for HTTP requests, we want this function to work for any TCP
// data proxying.
func proxyData(from net.Conn, to net.Conn, complete chan bool,
	done chan bool, otherDone chan bool) {
	var err error = nil
	var bytes []byte = make([]byte, 1024)
	var read int = 0
	for {
		select {
		// If the other channel is done processing, mark this channel
		// as done and move on
		case <-otherDone:
			complete <- true
			return
		default:
			// Enforce a small read deadline to make sure there is no
			// bottleneck
			from.SetReadDeadline(time.Now().Add(time.Second * 5))
			read, err = from.Read(bytes)
			// If any errors occured, write to complete as we are done (one of the
			// connections closed.)
			if err != nil {
				complete <- true
				done <- true
				logging.Error("Error while proxying request", err)
				return
			}
			// Write data to the destination.
			to.SetWriteDeadline(time.Now().Add(time.Second * 5))
			_, err = to.Write(bytes[:read])
			if err != nil {
				complete <- true
				done <- true
				return
			}
		}
	}
}

// HandleRequest implementation for Outbound handler.
// This function will handle sending the request from
// the client to the destination server
func (outbound *OutboundHandler) HandleRequest(request *socks5.Request) error {

	client := request.ClientConnection
	remote := request.OutboundConnection

	// defer func() {
	// 	request.State = RequestStateTerminating
	// 	// remote.Close()
	// }()

	complete := make(chan bool, 2)
	ch1 := make(chan bool, 1)
	ch2 := make(chan bool, 1)

	go proxyData(client, remote, complete, ch1, ch2)
	go proxyData(remote, client, complete, ch2, ch1)

	<-complete
	<-complete

	return nil

}
