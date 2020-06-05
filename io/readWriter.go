package io

import (
	"bufio"
	"fmt"
	"net"
)

// TunnelReadWriter is an abstraction to
// a read/writer used for proxying.
type TunnelReadWriter struct {
	connection     net.Conn
	peerConnection net.Conn

	// Inbound read/writers
	inReader *bufio.Reader
	inWriter *bufio.Writer

	// Outboud read/writers
	outReader *bufio.Reader
	outWriter *bufio.Writer
}

// NewTunnelReadWriter returns a new instance of TunnelReadWriter
func NewTunnelReadWriter(connection, peerConnection net.Conn) *TunnelReadWriter {
	inReader := bufio.NewReader(connection)
	inWriter := bufio.NewWriter(connection)

	outReader := bufio.NewReader(connection)
	outWriter := bufio.NewWriter(connection)

	rw := &TunnelReadWriter{
		connection:     connection,
		peerConnection: peerConnection,
		inReader:       inReader, inWriter: inWriter,
		outReader: outReader, outWriter: outWriter,
	}
	return rw
}

func newBufferReadWriter(inReader, outReader *bufio.Reader, inWriter, outWriter *bufio.Writer) *TunnelReadWriter {
	rw := &TunnelReadWriter{
		inReader: inReader, inWriter: inWriter,
		outReader: outReader, outWriter: outWriter,
	}

	return rw
}

func (rw *TunnelReadWriter) TunnelRead() error {
	var buffer []byte = make([]byte, 1024)
	read, err := rw.inReader.Read(buffer)
	fmt.Printf("Read : %d with val %s\n", read, string(buffer))
	if err != nil {
		fmt.Println("Error in read")
		return err
	}

	n, err := rw.outWriter.Write(buffer[:read])
	if err != nil {
		return err
	}
	rw.outWriter.Flush()
	fmt.Printf("Wrote : %d\n", n)
	return err
}
