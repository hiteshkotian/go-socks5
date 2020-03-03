package server

import (
	"fmt"
	"net"
	"testing"
	"time"
)

type testAddr struct {
	address string
}

func (a testAddr) Network() string {
	return "tcp"
}

func (a testAddr) String() string {
	return a.address
}

type testConn struct {
	msg    string
	result *string
}

func (c testConn) Read(b []byte) (n int, err error) {
	copy(b[:], c.msg)
	return len(c.msg), nil
}

func (c testConn) Write(b []byte) (n int, err error) {
	*c.result = string(b)
	return len(*c.result), nil
}

func (c testConn) Close() error {
	return nil
}

func (c testConn) LocalAddr() net.Addr {
	addr := testAddr{"128.0.0.1"}
	return addr
}

func (c testConn) RemoteAddr() net.Addr {
	addr := testAddr{"128.0.0.1"}
	return addr
}

func (c testConn) SetDeadline(t time.Time) error {
	return nil
}

func (c testConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c testConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestInvalidPort(t *testing.T) {
	s, err := NewServer(0)
	if s != nil {
		t.Errorf("Structure should not be returned for invalid port")
	}

	if err == nil {
		t.Errorf("Invalid error message provided")
	}
}

func TestInvalidPort2(t *testing.T) {
	s, err := NewServer(-100)
	if s != nil {
		t.Errorf("Structure should not be returned for invalid port")
	}

	if err == nil {
		t.Errorf("Invalid error message provided")
	}
}

func TestDefaultCreate(t *testing.T) {
	s, err := NewServer(8080)
	if err != nil {
		t.Errorf("Error returned when not expected")
	}

	if s == nil {
		t.Errorf("Invalid server instance returned")
	}
}

func TestHandleConnectionForHello(t *testing.T) {
	var result string
	c := &testConn{msg: "HELLO", result: &result}
	handleConnection(*c)

	fmt.Printf("result for %p is : %s\n", &c, result)

	if result != "Hello from 128.0.0.1\n" {
		t.Errorf("Invalid response : %s", result)
	}
}

func TestHandleConnectionForTime(t *testing.T) {
	var result string
	c := &testConn{msg: "TIME", result: &result}
	handleConnection(*c)

	fmt.Printf("result for %p is : %s\n", &c, result)
	// TODO See if we can check the exact response
	if len(result) < 1 {
		t.Errorf("Invalid response from the handler")
	}
}

func TestHandleConnectionHTTPGet(t *testing.T) {
	var result string
	expected := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nHost: localhost\r\n\r\nHello from the server.\nI got your message\r\n"
	c := &testConn{msg: "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n", result: &result}
	handleConnection(*c)

	fmt.Printf("Expecting : \n%s\n", expected)
	fmt.Printf("result for is : \n%s\n", result)
	if result != expected {
		t.Errorf("response string does not match. Expected : %s", expected)
	}
}

func TestHTTPGetMessage(t *testing.T) {
	msg := "GET / HTTP1/1\r\nHost: localhost\r\n\r\n"
	if !isHTTPMessage(msg) {
		t.Errorf("HTTP Get message test failed")
	}
}

func TestNonHTTPMessage(t *testing.T) {
	msg := "This is some random test message"
	if isHTTPMessage(msg) {
		t.Errorf("this is not an http test message")
	}
}
