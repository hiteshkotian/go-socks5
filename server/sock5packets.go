//File		: packets.go
//Author	: Nikhil Kotian
//Copyright	:

package server

import (
	"errors"
)

type nmethods uint8
type method uint8
type cmd uint8
type atype uint8
type reply uint8

const (
	//Socks5 Version field of the socks protocol
	Socks5 uint8 = 0x05
	//MaxMethodSize Maximum number of methods possible in a request packet
	MaxMethodSize nmethods = 0xFF

	//MethodNoAuth No authentication method
	MethodNoAuth method = 0x00
	//MethodGssAPI GSSAPI authentication method
	MethodGssAPI method = 0x01
	//MethodUserAuth User authenication method
	MethodUserAuth method = 0x02
	//MethodNoAcceptable No accpetable method found
	MethodNoAcceptable method = 0xFF

	//CmdConnect Client command for a connection
	CmdConnect cmd = 0x01
	//CmdBind Client command for a bind
	CmdBind cmd = 0x02
	//CmdUDPAssc Client command for UDP Asscociate
	CmdUDPAssc cmd = 0x03

	//AtypIPV4 IPV4 Address type
	AtypIPV4 atype = 0x01
	//AtypDomain Domain Address type
	AtypDomain atype = 0x03
	//AtypIPV6 IPV6 Address type
	AtypIPV6 atype = 0x04

	//AddrIPV4Size IPV4 Address size
	AddrIPV4Size uint8 = 0x04
	//AddrIPV6Size IPV6 Address size
	AddrIPV6Size uint8 = 0x10

	//ReplySucceeded Reply sent to client suceeded
	ReplySucceeded reply = 0x00
	//ReplyGeneralFail Reply sent to client that the request failed
	ReplyGeneralFail reply = 0x01
	//ReplyConnDenied Reply sent to client that connection is denied
	ReplyConnDenied reply = 0x02
	//ReplyNetUnreachable Reply sent to client that network is unreachable
	ReplyNetUnreachable reply = 0x03
	//ReplyHostUnreachable Reply sent to client that the host is unreachable
	ReplyHostUnreachable reply = 0x04
	//ReplyConnRefused Reply sent to client that connection was refused
	ReplyConnRefused reply = 0x05
	//ReplyTTLExpired Reply sent to client that TTL expired
	ReplyTTLExpired reply = 0x06
	//ReplyCmdUnsupp Reply sent to client that command is unsupported
	ReplyCmdUnsupp reply = 0x07
	//ReplyAddrTypUnsupp Reply sent to client that address type is unsupported
	ReplyAddrTypUnsupp reply = 0x08
)

//MethodSelectionReq Method selection request packet
type MethodSelectionReq struct {
	nmethods
	methods []method
}

//MethodSelectionResp Method selection response packet
type MethodSelectionResp struct {
	method
}

//SockRequest Sock5 Request structure
type SockRequest struct {
	cmd
	atype
	destaddr []uint8
	destport uint16
}

//SockReply reply structure
type SockReply struct {
	reply
	atype
	bindaddr []uint8
	bindport uint16
}

func checkMessageVersion(msg []uint8) error {
	if msg[0] != Socks5 {
		return errors.New("Sock5Packet: Version incorrect")
	}
	return nil
}

//GetSocketMethod Used to decode Method packet
func GetSocketMethod(msg []uint8) (MethodSelectionReq, error) {
	var ret MethodSelectionReq

	if len(msg) < 2 || len(msg) != int(msg[1])+2 {
		return ret, errors.New("Sock5Packet: Message incorrect size in method negotiation")
	}

	if err := checkMessageVersion(msg); err != nil {
		return ret, err
	}

	if len(msg) > int(MaxMethodSize) {
		return ret, errors.New("Sock5Packet: Message size too big in method negotation")
	}

	ret.nmethods = nmethods(msg[1])

	for _, v := range msg[2:] {
		ret.methods = append(ret.methods, method(v))
	}
	return ret, nil
}

//GetSocketMethodResponse Get serialized Method response
func GetSocketMethodResponse(resp MethodSelectionResp) ([]uint8, error) {
	ret := make([]uint8, 2)
	ret[0] = uint8(Socks5)
	ret[1] = uint8(resp.method)
	return ret, nil
}

//GetSocketRequestDeserialized Get deserialized socket request
func GetSocketRequestDeserialized(msg []uint8) (SockRequest, error) {
	var ret SockRequest

	if len(msg) < 5 {
		return ret, errors.New("Socks5Packet: Packet too small for a request")
	}

	if err := checkMessageVersion(msg); err != nil {
		return ret, err
	}

	ret.cmd = cmd(msg[1])
	ret.atype = atype(msg[3])

	var size uint8
	var addrStart uint8 = 4

	switch ret.atype {
	case AtypIPV4:
		size = AddrIPV4Size
	case AtypIPV6:
		size = AddrIPV6Size
	case AtypDomain:
		size = msg[4]
		addrStart = 5
	default:
		return ret, errors.New("Socks5Packet: Wrong address type in request")
	}

	if len(msg) > 255 || addrStart+size+2 != uint8(len(msg)) {
		return ret, errors.New("Socks5Packet: Packet wrong size in request")
	}

	ret.destaddr = make([]uint8, size)
	counter := 0
	for _, v := range msg[addrStart : addrStart+size] {
		ret.destaddr[counter] = v
		counter++
	}

	ret.destport = (uint16(msg[addrStart+size+1]) << 8) | uint16(msg[addrStart+size])
	return ret, nil
}

//GetSocketResponseSerialized Serializes the socket reply
func GetSocketResponseSerialized(resp SockReply) ([]uint8, error) {
	ret := make([]uint8, 4)

	ret[0] = Socks5
	ret[1] = uint8(resp.reply)
	ret[2] = 0x00
	ret[3] = uint8(resp.atype)

	var size uint8

	switch resp.atype {
	case AtypIPV4:
		size = AddrIPV4Size
	case AtypIPV6:
		size = AddrIPV6Size
	case AtypDomain:
		size = uint8(len(resp.bindaddr))
		ret = append(ret, size)
	default:
		return ret, errors.New("Socks5Packet: Wrong address type in response")
	}

	if len(resp.bindaddr) != int(size) {
		return ret, errors.New("Socks5Packet: Bind address size is not same as type")
	}

	for _, v := range resp.bindaddr {
		ret = append(ret, v)
	}

	ret = append(ret, uint8(resp.bindport&0xFF))
	ret = append(ret, uint8((resp.bindport>>8)&0xFF))

	return ret, nil
}
