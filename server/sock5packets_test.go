package server

import (
	"testing"
)

func CompareSlices(a, b []uint8) bool{
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
	msg := []uint8{1,2,3}
	_, err := GetSocketMethod(msg)
	if err == nil {
		t.Errorf("Sock5Packet: Method Selection for Socks 5 not working")
	}
}

func TestMethodSelectionCorrect(t *testing.T) {
	msg := []uint8{Socks5, 0x2, uint8(MethodNoAuth), uint8(MethodGssAPI)}
	methodMsg, err := GetSocketMethod(msg)
	
	if err != nil {
		t.Errorf("Sock5Packet: Method selection for Socks5 not working")
	}

	if methodMsg.nmethods != nmethods(2) {
		t.Errorf("Sock5Packet: Request method selection number of methods incorrect")
	}

	if len(methodMsg.methods) != 2 {
		t.Errorf("Sock5Packet: Request wrong methods size")
	}

	if methodMsg.methods[0] != MethodNoAuth && methodMsg.methods[1] != MethodGssAPI {
		t.Errorf("Sock5Packet: Request wrong methods")
	} 
}

func TestMethodSelectionSizeBigger(t *testing.T) {
	msg := make([]uint8, 512)
	_, err := GetSocketMethod(msg)
	if err == nil {
		t.Errorf("Sock5Packet: Bigger selection size not flagged as error")
	}
}

func TestWrongSize(t *testing.T) {
	msg := []uint8 {Socks5, 0x2, uint8(MethodGssAPI)}
	_, err := GetSocketMethod(msg)
	if err == nil {
		t.Errorf("Sock5Packet: Wrong size not detected")
	}

	shortmsg := []uint8 {uint8(Socks5)}
	_, err = GetSocketMethod(shortmsg)
	if err == nil {
		t.Errorf("Sock5Packet: Short size not detected")
	}
}

func TestResponseEncoding(t *testing.T) {
	respMsg := MethodSelectionResp {MethodUserAuth}
	resp, err := GetSocketMethodResponse(respMsg)
	if err != nil {
		t.Errorf("Sock5Packet: Method Response does not work")
	}

	if resp[0] != Socks5 && resp[1] != uint8(MethodUserAuth) {
		t.Errorf("Sock5Packet: Method Response encoding wrong")
	}
}

func TestSockRequestDecodeWrongVersion(t *testing.T) {
	msg := []uint8 {0x2,0x3,0x4}
	_, err := GetSocketRequestDeserialized(msg)
	if err == nil {
		t.Errorf("Sock5Packet: Socket Request does not work")
	}
}

func TestSockRequestDecodeCorrectIPV4(t *testing.T) {
	msg := []uint8 {Socks5, uint8(CmdConnect), 0x00, uint8(AtypIPV4), 0x1, 0x4, 0x5, 0x21, 0x45, 0x56}
	resp, err := GetSocketRequestDeserialized(msg)
	
	if err != nil {
		t.Errorf("Sock5Packet: Socket Request Deserialization failed with IPV4")
	}

	if resp.cmd != CmdConnect || resp.atype != AtypIPV4 {
		t.Errorf("Sock5Packet: Socket Request Deserialization failed with address data")
	} 
	
	if len(resp.destaddr) != 4 {
		t.Errorf("Sock5Packet: IPV4 address wrong size")
	}

	if resp.destaddr[0] != 0x1 || resp.destaddr[1] != 0x4 || resp.destaddr[2] != 0x5 || resp.destaddr[3] != 0x21{
		t.Errorf("Sock5Packet: IPV4 address wrong")
	}

	if resp.destport != 0x5645 {
		t.Errorf("Sock5Packet: IPV4 port wrong")
	}
}

func TestSockRequestDecodeCorrectIPV6(t *testing.T) {
	msg := []uint8 {Socks5, uint8(CmdConnect), 0x00, uint8(AtypIPV6), 0x1, 0x4, 0x5, 0x21,
																			  0x22, 0x24, 0x36, 0x38,
																			  0x39, 0x41, 0x43, 0x44, 
																			  0x45, 0x46, 0x47, 0x48,
																			  0x45, 0x56}
	resp, err := GetSocketRequestDeserialized(msg)
	
	if err != nil {
		t.Errorf("Sock5Packet: Socket Request Deserialization failed with IPV4")
	}

	if resp.cmd != CmdConnect || resp.atype != AtypIPV6 {
		t.Errorf("Sock5Packet: Socket Request Deserialization failed with address data")
	} 
	
	if len(resp.destaddr) != 16 {
		t.Errorf("Sock5Packet: IPV6 address wrong size")
	}

	counter := 0
	for _,v := range msg[4:19] {
		if resp.destaddr[counter] != v {
			t.Errorf("Sock5Packet: IPV6 address wrong")
			return
		}
		counter++
	}

	if resp.destport != 0x5645 {
		t.Errorf("Sock5Packet: IPV6 port wrong")
	}
}

func TestSockRequestDecodeCorrectDomain(t *testing.T) {
	msg := []uint8 {Socks5, uint8(CmdConnect), 0x00, uint8(AtypDomain), 0x3, 0xA, 0x5, 0x21,
																			  0x45, 0x56}
	resp, err := GetSocketRequestDeserialized(msg)
	
	if err != nil {
		t.Errorf("Sock5Packet: Socket Request Deserialization failed with DomainType")
	}

	if resp.cmd != CmdConnect || resp.atype != AtypDomain {
		t.Errorf("Sock5Packet: Socket Request Deserialization failed with address data")
	} 
	
	if len(resp.destaddr) != 3 {
		t.Errorf("Sock5Packet: Domain address wrong size")
	}

	counter := 0
	for _,v := range msg[5:8] {
		if resp.destaddr[counter] != v {
			t.Errorf("Sock5Packet: Domain address wrong")
			return
		}
		counter++
	}

	if resp.destport != 0x5645 {
		t.Errorf("Sock5Packet: Domain port wrong")
	}
}

func TestSocketRequestDecodeWrongSize(t *testing.T) {
	msg := []uint8 {Socks5, uint8(CmdConnect)}

	_, err := GetSocketRequestDeserialized(msg)

	if err == nil {
		t.Errorf("Sock5Packet: Didnot detect smaller size")
	}

	wrongsize := []uint8 {Socks5, uint8(CmdConnect), 0x00, uint8(AtypIPV4), 0x1, 0x4, 0x21, 0x45, 0x56}
	_, rerr := GetSocketRequestDeserialized(wrongsize)

	if rerr == nil {
		t.Errorf("Sock5Packet: Didnot detect wrong size")
	}
}

func TestSocketResponseCorrect(t *testing.T) {
	resp := SockReply {ReplyGeneralFail, AtypIPV4, []uint8{0x5,0x6,0x7,0x8}, 0x5645}

	msg, _ := GetSocketResponseSerialized(resp)
	
	if msg[0] != Socks5 || msg[1] != uint8(ReplyGeneralFail) || msg[2] != 0x00 || msg[3] != uint8(AtypIPV4)	{
		t.Errorf("Sock5Packet: Version is wrong")
	}

	counter := 0
	for _, v := range msg[4:7] {
		if v != resp.bindaddr[counter] {
			t.Errorf("Sock5Packet: IPV4 address in reply does not match")
		}
		counter++
	}

	if msg[8] != 0x45 || msg[9] != 0x56 {
		t.Errorf("Sock5Packet: Bind port is wrong in reply")
	}
}

func TestSocketResponseIPV6(t *testing.T) {
	resp := SockReply {ReplyGeneralFail, AtypIPV6, []uint8{0x3, 0x4, 0x5, 0x6, 
														   0x7, 0x8, 0x9, 0xa,
														   0xb, 0xc, 0xd, 0xe,
														   0xf, 0x10, 0x11, 0x12}, 0x5467}

	msg, _ := GetSocketResponseSerialized(resp)
	
	if msg[0] != Socks5 || msg[1] != uint8(ReplyGeneralFail) || msg[2] != 0x00 || msg[3] != uint8(AtypIPV6){
		t.Errorf("Socks5Packet: Socks reply IPV6 invalid")
	}

	counter := 0
	for _, v := range msg[4:19] {
		if resp.bindaddr[counter] != v {
			t.Errorf("Socks5Packet: IPV6 address wrongly encoded")
		}
		counter++
	}

	if msg[20] != 0x67 || msg[21] != 0x54 {
		t.Errorf("Socks5Packet: Bind port is incorrect")
	}
}

func TestSocketResponseDomain(t *testing.T) {
	resp := SockReply {ReplyGeneralFail, AtypDomain, []uint8{0x34, 0x35, 0x36}, 0x5434}

	msg, _ := GetSocketResponseSerialized(resp)

	if msg[0] != Socks5 || msg[1] != uint8(ReplyGeneralFail) || msg[2] != 0x00 || msg[3] != uint8(AtypDomain) || msg[4] != 3{
		t.Errorf("Socks5Packet: Socks reply domain invalid")
	}

	if msg[5] != 0x34 || msg[6] != 0x35 || msg[7] != 0x36 {
		t.Errorf("Socks5Packet: Domain address decoding response is incorrect")
	} 
}

func TestSocketResponseIncorrectSize(t *testing.T) {
	resp := SockReply {ReplyGeneralFail, AtypIPV4, []uint8{0x34, 0x35}, 0x5434}

	_,err := GetSocketResponseSerialized(resp)

	if err == nil {
		t.Errorf("Socks5Packet: Incorrect size not caught by serialization")
	}
}