package proxy

import (
	"hiteshkotian/ssl-tunnel/logging"
	"hiteshkotian/ssl-tunnel/socks5"
	"net"
	"testing"
)

// Test if socks5 init method works as expected
func TestInitWithNoMethods(t *testing.T) {
	var clientRequest []byte

	expectedVersion := uint8(0x05)
	// Version
	clientRequest = append(clientRequest, expectedVersion)
	// Set Methods
	clientRequest = append(clientRequest, 0x00)
	clientRequest = append(clientRequest, 0x00)

	server := Server{name: "test"}
	writeConn, readConn := net.Pipe()
	defer readConn.Close()

	go func() {
		request := &socks5.Request{ClientConnection: writeConn}
		server.handleInitialLocal(request)
		writeConn.Close()
	}()

	readConn.Write(clientRequest)

	response := make([]byte, 512)
	n, e := readConn.Read(response)
	if e != nil {
		t.Error("Error reading response: ", e)
	} else if n <= 0 {
		t.Error("Invalid bytes read")
	}

	// check version
	if response[0] != expectedVersion {
		t.Errorf("Invalid version. Expected 0x%02x, received 0x%0x2",
			expectedVersion, response[0])
	}

	if response[1] != 0x00 {
		t.Errorf("Invalid command selected. Expected 0x00, received : 0x%02x", response[1])
	}
}

func TestInitWithInvalidVersion(t *testing.T) {
	var clientRequest []byte

	expectedVersion := uint8(0x05)
	expectedCMD := uint8(0xff)

	// Version
	clientRequest = append(clientRequest, 0x0a)
	// Set Methods
	clientRequest = append(clientRequest, 0x00)

	server := Server{name: "test"}
	writeConn, readConn := net.Pipe()

	defer readConn.Close()

	go func() {
		request := &socks5.Request{ClientConnection: writeConn}
		server.handleInitialLocal(request)
		writeConn.Close()
	}()

	readConn.Write(clientRequest)

	response := make([]byte, 512)
	n, e := readConn.Read(response)
	if e != nil {
		t.Error("Error reading response: ", e)
	} else if n <= 0 {
		t.Error("Invalid bytes read")
	}

	logging.DumpHex(response[:n], "Response is")

	if response[0] != expectedVersion {
		t.Errorf("Invalid version. Expected 0x%02x, received 0x%0x2",
			expectedVersion, response[0])
	}

	if response[1] != expectedCMD {
		t.Errorf("Invalid command sent. Expected 0x%02x, received 0x%0x2",
			expectedCMD, response[1])
	}
}

func TestInitWithIncompletePacket(t *testing.T) {
	var clientRequest []byte

	expectedVersion := uint8(0x05)
	// Version
	clientRequest = append(clientRequest, expectedVersion)
	// Set Methods
	clientRequest = append(clientRequest, 0x00)

	server := Server{name: "test"}
	writeConn, readConn := net.Pipe()
	defer readConn.Close()

	go func() {
		request := &socks5.Request{ClientConnection: writeConn}
		server.handleInitialLocal(request)
		writeConn.Close()
	}()

	readConn.Write(clientRequest)

	response := make([]byte, 512)
	n, e := readConn.Read(response)
	if e != nil {
		t.Error("Error reading response: ", e)
	} else if n <= 0 {
		t.Error("Invalid bytes read")
	}

	// check version
	if response[0] != expectedVersion {
		t.Errorf("Invalid version. Expected 0x%02x, received 0x%0x2",
			expectedVersion, response[0])
	}

	if response[1] != 0xff {
		t.Errorf("Invalid command selected. Expected 0xff, received : 0x%02x", response[1])
	}
}
