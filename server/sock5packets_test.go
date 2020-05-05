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

func TestUDPHeaderIPV4Request(t *testing.T){
	resp := []uint8 {0x00, 0x00, 0x54, uint8(AtypIPV4), 0x34, 0x43, 0x55, 0x56, 0x23, 0x45, 0x32, 0x34, 0x45}

	respDeserial, err := GetSocketUDPDeserialized(resp)

	if err != nil {
		t.Errorf("Socks5Packet: Decoding of UDP packet failed")
	}

	if respDeserial.fragment != 0x54 || respDeserial.atype != AtypIPV4 {
		t.Errorf("Socks5Packet: Incorrect decoding of UDP packet")
	}

	for i,v := range resp[4:8] {
		if v != respDeserial.address[i] {
			t.Errorf("Socks5Packet: Incorrect decoding of UDP packet IPV4 address")
		}
	}

	if respDeserial.port != 0x4523 {
		t.Errorf("Socks5Packet: Inccorect decoding of UDP packet IPV4 port")
	}

	for i,v := range resp[10:] {
		if v != respDeserial.data[i] {
			t.Errorf("Socks5Packet: Incorrect decoding of UDP packet IPV4 data")
		}
	}
}

func TestUDPHeaderIPV6Request(t *testing.T) {
	resp := []uint8 {0x00, 0x00, 0x95, uint8(AtypIPV6), 0x32, 0x45, 0x56, 0x65, 
														0x56, 0x54, 0x65, 0x56,
														0x55, 0x54, 0x78, 0x9a,
														0x59, 0x78, 0xbb, 0xee, 
														0x45, 0x55,
														0x55, 0xbe, 0xab, 0xcd}
	
	respDeserial, err := GetSocketUDPDeserialized(resp)
	if err != nil {
		t.Errorf("Socks5Packet: Decoding of UDP packet failed")
	}

	if respDeserial.fragment != 0x95 || respDeserial.atype != AtypIPV6 {
		t.Errorf("Socks5Packet: Incorrect decoding of UDP packet")
	}	

	for i,v := range resp[4:20] {
		if v != respDeserial.address[i] {
			t.Errorf("Socks5Packet: Incorrect decoding of UDP packet IPV6 address")
		}
	}	

	if respDeserial.port != 0x5545 {
		t.Errorf("Socks5Packet: Inccorect decoding of UDP packet IPV6 port")
	}

	for i,v := range resp[22:] {
		if v != respDeserial.data[i] {
			t.Errorf("Socks5Packet: Incorrect decoding of UDP packet IPV4 data")
		}
	}
}

func TestUDPHeaderIPV6DomainName(t *testing.T) {
	resp := []uint8 {0x00, 0x00, 0x99, uint8(AtypDomain), 0x13, 0x45, 0x56, 0x65, 
														0x56, 0x54, 0x65, 0x56,
														0x55, 0x54, 0x78, 0x9a,
														0x59, 0x78, 0xbb, 0xee, 
														0x45, 0x55, 0x67, 0x99,
														0x55, 0xbe, 0xab, 0xcd}
	
	respDeserial, err := GetSocketUDPDeserialized(resp)
	if err != nil {
		t.Errorf("Socks5Packet: Decoding of UDP packet failed")
	}

	if respDeserial.fragment != 0x99 || respDeserial.atype != AtypDomain {
		t.Errorf("Socks5Packet: Incorrect decoding of UDP packet")
	}

	for i,v := range resp[5:24] {
		if v != respDeserial.address[i] {
			t.Errorf("Socks5Packet: Incorrect decoding of UDP packet IPV6 address")
		}
	}

	if respDeserial.port != 0xbe55 || respDeserial.data[0] != 0xab || respDeserial.data[1] != 0xcd {
		t.Errorf("Socks5Packet: Incorrect decoding of UDP packet IPV6 port")
	}
}

func TestUDPHeaderSmallSize(t *testing.T) {
	resp := []uint8 {0x00, 0x00, 0x99}
	
	_, err := GetSocketUDPDeserialized(resp)

	if err == nil {
		t.Errorf("Sock5Packet: Incorrect decoding not detected")
	}
}

func TestUDPSerializeIPV4(t *testing.T) {
	resp := UDPPacket{fragment: 0x45, atype: AtypIPV4, address: []uint8{0x34,0x45,0x56,0x78}, 
						port: 0x4556, data: []uint8{0x34, 0x98, 0x78, 0x99}}
   
	ret, err := GetSocketUDPSerialized(resp)
	
	if err != nil {
		t.Errorf("Sock5Packet: Error serializing IPV4 packet")
	}

	if ret[0] != 0x00 || ret[1] != 0x00 || ret[2] != 0x45 || ret[3] != uint8(AtypIPV4){
		t.Errorf("Sock5Packet: Error serializing IPV4 packet header")
	}

	if ret[4] != 0x34 || ret[5] != 0x45 || ret[6] != 0x56 || ret[7] != 0x78 {
		t.Errorf("Socks5Packet: Error serializing IPV4 packet address")
	}

	if ret[8] != 0x56 || ret[9] != 0x45 {
		t.Errorf("Sock5Packet: Error serializin IPV4 packet port")
	}

	if ret[10] != 0x34 || ret[11] != 0x98 || ret[12] != 0x78 || ret[13] != 0x99 {
		t.Errorf("Socks5Packet: Error serializing IPV4 packet data")
	}
}

func TestUDPPacketSerializeIPV6(t *testing.T) {
	resp := UDPPacket{fragment : 0x78, atype: AtypIPV6, address: []uint8{0x78, 0x67, 0x88, 0x89, 
																		0x87, 0x76, 0x88, 0x98,
																		0x77, 0x66, 0x87, 0x88,
																		0x77, 0x66, 0x78, 0x88},
														port: 0x6734, 
														data: []uint8{0x78, 0x99}}

	ret, err := GetSocketUDPSerialized(resp)

	if err != nil {
		t.Errorf("Sock5Packet: Error serializing IPV6 packet")
	}

	if ret[0] != 0x00 || ret[1] != 0x00 || ret[2] != 0x78 || ret[3] != uint8(AtypIPV6){
		t.Errorf("Sock5Packet: Error serializing IPV6 packet header")
	}

	for i,v := range(ret[4:20]) {
		if resp.address[i] != v {
			t.Errorf("Sock5Packt: Error serializing IPV6 packet address")
		}
	}

	if ret[20] != 0x34 || ret[21] != 0x67 {
		t.Errorf("Sock5Packet: Error serializin IPV6 packet port")
	}

	if ret[22] != 0x78 || ret[23] != 0x99  {
		t.Errorf("Socks5Packet: Error serializing IPV6 packet data")
	}
}

func TestUDPPacketSerializeDomain(t *testing.T) {
	resp := UDPPacket {fragment: 0xab, atype: AtypDomain, address: []uint8{0x1, 0x34, 0x45}, 
						port:0x4563, data: []uint8{0x7}}

	ret, err := GetSocketUDPSerialized(resp)

	if err != nil {
		t.Errorf("Sock5Packet: Error serializing domain packet")
	}
					
	if ret[0] != 0x00 || ret[1] != 0x00 || ret[2] != 0xab || ret[3] != uint8(AtypDomain) || ret[4] != 3{
		t.Errorf("Sock5Packet: Error serializing domain packet header")
	}
	
	for i,v := range(ret[5:8]) {
		if resp.address[i] != v {
			t.Errorf("Sock5Packt: Error serializing domain packet address")
		}
	}

	if ret[8] != 0x63 || ret[9] != 0x45 {
		t.Errorf("Sock5Packet: Error serializin IPV6 packet port")
	}

	if ret[10] != 0x7 {
		t.Errorf("Socks5Packet: Error serializing IPV6 packet data")
	}
}