package server

import (
	"testing"
)

func CompareSlices(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestMethodSelectionDecodeCheck(t *testing.T) {
	msg := []uint8{1, 2, 3}
	_, err := GetSocketMethod(msg)
	if err == nil {
		t.Errorf("Sock5: Method Selection for Socks 5 not working")
	}
}

func TestMethodSelectionCorrect(t *testing.T) {
	msg := []uint8{uint8(Socks5), 0x2, uint8(MethodNoAuth), uint8(MethodGssAPI)}
	methodMsg, err := GetSocketMethod(msg)

	if err != nil {
		t.Errorf("Sock5: Method selection for Socks5 not working")
	}

	if methodMsg.ver != Socks5 {
		t.Errorf("Sock5: Request method selection Socket version error")
	}

	if methodMsg.nmethods != nmethods(2) {
		t.Errorf("Sock5: Request method selection number of methods incorrect")
	}

	if len(methodMsg.methods) != 2 {
		t.Errorf("Socks5: Request wrong methods size")
	}

	if methodMsg.methods[0] != MethodNoAuth && methodMsg.methods[1] != MethodGssAPI {
		t.Errorf("Socks5: Request wrong methods")
	}
}

func TestMethodSelectionSizeBigger(t *testing.T) {
	msg := make([]uint8, 512)
	_, err := GetSocketMethod(msg)
	if err == nil {
		t.Errorf("Sock5: Bigger selection size not flagged as error")
	}
}

func TestWrongSize(t *testing.T) {
	msg := []uint8{uint8(Socks5), 0x2, uint8(MethodGssAPI)}
	_, err := GetSocketMethod(msg)
	if err == nil {
		t.Errorf("Sock5: Wrong size not detected")
	}

	shortmsg := []uint8{uint8(Socks5)}
	_, err = GetSocketMethod(shortmsg)
	if err == nil {
		t.Errorf("Sock5:Short size not detected")
	}
}

func TestResponseEncoding(t *testing.T) {
	respMsg := MethodSelectionResp{Socks5, MethodUserAuth}
	resp, err := GetSocketMethodResponse(respMsg)
	if err != nil {
		t.Errorf("Socks: Method Response does not work")
	}

	if resp[0] != uint8(Socks5) && resp[1] != uint8(MethodUserAuth) {
		t.Errorf("Socks: Method Response encoding wrong")
	}
}

func TestSockRequestDecodeWrongVersion(t *testing.T) {
	msg := []uint8{0x2, 0x3, 0x4}
	_, err := GetSocketRequestDeserialized(msg)
	if err == nil {
		t.Errorf("Socks: Socket Request does not work")
	}
}

func TestSockRequestDecodeCorrectIPV4(t *testing.T) {
	msg := []uint8{uint8(Socks5), uint8(CmdConnect), 0x00, uint8(AtypIPV4), 0x1, 0x4, 0x5, 0x21, 0x45, 0x56}
	resp, err := GetSocketRequestDeserialized(msg)

	if err != nil {
		t.Errorf("Socks: Socket Request Deserialization failed with IPV4")
	}

	if resp.ver != Socks5 || resp.cmd != CmdConnect || resp.atype != AtypIPV4 {
		t.Errorf("Socks: Socket Request Deserialization failed with address data")
	}

	if len(resp.destaddr) != 4 {
		t.Errorf("Socks: IPV4 address wrong size")
	}

	if resp.destaddr[0] != 0x1 || resp.destaddr[1] != 0x4 || resp.destaddr[2] != 0x5 || resp.destaddr[3] != 0x21 {
		t.Errorf("Socks: IPV4 address wrong")
	}

	if resp.destport != 0x5645 {
		t.Errorf("Socks: IPV4 port wrong")
	}
}

func TestSockRequestDecodeCorrectIPV6(t *testing.T) {
	msg := []uint8{uint8(Socks5), uint8(CmdConnect), 0x00, uint8(AtypIPV6), 0x1, 0x4, 0x5, 0x21,
		0x22, 0x24, 0x36, 0x38,
		0x39, 0x41, 0x43, 0x44,
		0x45, 0x46, 0x47, 0x48,
		0x45, 0x56}
	resp, err := GetSocketRequestDeserialized(msg)

	if err != nil {
		t.Errorf("Socks: Socket Request Deserialization failed with IPV4")
	}

	if resp.ver != Socks5 || resp.cmd != CmdConnect || resp.atype != AtypIPV6 {
		t.Errorf("Socks: Socket Request Deserialization failed with address data")
	}

	if len(resp.destaddr) != 16 {
		t.Errorf("Socks: IPV6 address wrong size")
	}

	counter := 0
	for _, v := range msg[4:19] {
		if resp.destaddr[counter] != v {
			t.Errorf("Socks: IPV6 address wrong")
			return
		}
		counter++
	}

	if resp.destport != 0x5645 {
		t.Errorf("Socks: IPV6 port wrong")
	}
}

func TestSockRequestDecodeCorrectDomain(t *testing.T) {
	msg := []uint8{uint8(Socks5), uint8(CmdConnect), 0x00, uint8(AtypDomain), 0x3, 0xA, 0x5, 0x21,
		0x45, 0x56}
	resp, err := GetSocketRequestDeserialized(msg)

	if err != nil {
		t.Errorf("Socks: Socket Request Deserialization failed with DomainType")
	}

	if resp.ver != Socks5 || resp.cmd != CmdConnect || resp.atype != AtypDomain {
		t.Errorf("Socks: Socket Request Deserialization failed with address data")
	}

	if len(resp.destaddr) != 3 {
		t.Errorf("Socks: Domain address wrong size")
	}

	counter := 0
	for _, v := range msg[5:8] {
		if resp.destaddr[counter] != v {
			t.Errorf("Socks: Domain address wrong")
			return
		}
		counter++
	}

	if resp.destport != 0x5645 {
		t.Errorf("Socks: Domain port wrong")
	}
}
