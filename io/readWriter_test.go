package io

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestTunnel(t *testing.T) {
	s := "Hello"
	out := ""
	buf := bytes.NewBufferString(s)

	outBuf := bytes.NewBufferString(out)

	inReader := bufio.NewReader(buf)
	inWriter := bufio.NewWriter(buf)

	outReader := bufio.NewReader(outBuf)
	outWriter := bufio.NewWriter(outBuf)

	rw := newBufferReadWriter(inReader, outReader, inWriter, outWriter)
	fmt.Fprintf(inWriter, "There there")

	err := rw.TunnelRead()
	if err != nil {
		t.Errorf("Error on tunnel write : %s", err.Error())
	}

	outBuff := make([]byte, 20)
	n, err := outReader.Read(outBuff)

	output := string(outBuff)
	fmt.Printf("output is : %s [%d]\n", output, n)
}
