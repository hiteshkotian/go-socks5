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
		t.Errorf("Sock5:Method Selection for Socks 5 not working")
	}
}

func TestMethodSelectionCorrect(t *testing.T) {
	msg := []uint8{uint8(Socks5), 0x2, uint8(MethodNoAuth), uint8(MethodGssAPI)}
	methodMsg, err := GetSocketMethod(msg)
	
	if err != nil {
		t.Errorf("Sock5:Method selection for Socks5 not working")
	}

	if methodMsg.ver != Socks5 {
		t.Errorf("Sock5:Request method selection Socket version error")
	}

	if methodMsg.nmethods != nmethods(2) {
		t.Errorf("Sock5:Request method selection number of methods incorrect")
	}

	if len(methodMsg.methods) != 2 {
		t.Errorf("Socks5:Request wrong methods size")
	}

	if methodMsg.methods[0] != MethodNoAuth && methodMsg.methods[1] != MethodGssAPI {
		t.Errorf("Socks5:Request wrong methods")
	} 
}

func TestMethodSelectionSizeBigger(t *testing.T) {
	msg := make([]uint8, 512)
	_, err := GetSocketMethod(msg)
	if err == nil {
		t.Errorf("Sock5:Bigger selection size not flagged as error")
	}
}

func TestWrongSize(t *testing.T) {
	msg := []uint8 {uint8(Socks5), 0x2, uint8(MethodGssAPI)}
	_, err := GetSocketMethod(msg)
	if err == nil {
		t.Errorf("Sock5:Wrong size not detected")
	}

	shortmsg := []uint8 {uint8(Socks5)}
	_, err = GetSocketMethod(shortmsg)
	if err == nil {
		t.Errorf("Sock5:Short size not detected")
	}
}