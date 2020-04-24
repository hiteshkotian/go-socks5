//File		: packets.go
//Author	: Nikhil Kotian
//Copyright	:

package server

type ver 			uint8
type nmethods 		uint8
type method 		uint8
type cmd			uint8
type rsv			uint8
type addresstype 	uint8
type atype			uint8
type addr		  []uint8
type port			uint16
type reply			uint8

const (
	SOCKS5 			ver		= 0x05
	MAXMETHODSIZE	nmethods	= 0xFF	

	METHOD_NOAUTH			method = 0x00
	METHOD_GSSAPI			method = 0x01
	METHOD_USERAUTH			method = 0x02
	METHOD_NOACCEPTABLE		method = 0xFF

	CMD_CONNECT			cmd	= 0x01
	CMD_BIND			cmd = 0x02
	CMD_UDPASSC			cmd = 0x03 	

	ATYP_IPV4			atype = 0x01
	ATYP_DOMAIN			atype = 0x03
	ATYP_IPV6			atype = 0x04

	ADDR_IPV4_SIZE		uint8 = 0x04
	ADDR_IPV6_SIZE		uint8 = 0x10

	REPLY_SUCCEEDED			reply = 0x00
	REPLY_GENERAL_FAIL		reply = 0x01
	REPLY_CONN_DENIED		reply = 0x02
	REPLY_NET_UNREACHABLE	reply = 0x03
	REPLY_HOST_UNREACHABLE	reply = 0x04
	REPLY_CONN_REFUSED		reply = 0x05
	REPLY_TTL_EXPIRED		reply = 0x06
	REPLY_CMD_UNSUPP		reply = 0x07
	REPLY_ADDRTYP_UNSUPP	reply = 0x08
)

type MethodSelectionReq struct {
	ver
	nmethods
	methods		[]method
}

type MethodSelectionResp struct {
	ver
	method
}

type SockRequest struct {
	ver
	cmd
	rsv
	addresstype
	destaddr		addr
	destport		port	
}

type SockReply struct {
	ver
	reply
	reserved uint8
	bindaddr		addr
	bindport		port
}