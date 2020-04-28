//File		: packets.go
//Author	: Nikhil Kotian
//Copyright	:

package server

import (
	"errors"
)

type ver uint8
type nmethods uint8
type method uint8
type cmd uint8
type rsv uint8
type addresstype uint8
type atype uint8
type addr []uint8
type port uint16
type reply uint8

const (
	//Socks5 Version field of the socks protocol
	Socks5 ver = 0x05
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
	ver
	nmethods
	methods []method
}

//MethodSelectionResp Method selection response packet
type MethodSelectionResp struct {
	ver
	method
}

//SockRequest Sock5 Request structure
type SockRequest struct {
	ver
	cmd
	rsv
	addresstype
	destaddr addr
	destport port
}

//SockReply reply structure
type SockReply struct {
	ver
	reply
	reserved uint8
	bindaddr addr
	bindport port
}

//GetSocketMethod Used to decode Method packet
func GetSocketMethod(msg []uint8) (MethodSelectionReq, error) {
	var ret MethodSelectionReq

	if len(msg) < 2 || len(msg) != int(msg[1])+2 {
		return ret, errors.New("Socks5:SOCKS Message incorrect size")
	}

	if msg[0] != uint8(Socks5) {
		return ret, errors.New("Sock5:SOCKS VERSION INCORRECT")
	}

	if len(msg) > int(MaxMethodSize) {
		return ret, errors.New("Sock5:SOCKS Message size too big")
	}

	ret.ver = ver(msg[0])
	ret.nmethods = nmethods(msg[1])

	for _, v := range msg[2:] {
		ret.methods = append(ret.methods, method(v))
	}
	return ret, nil
}
