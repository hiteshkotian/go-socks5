package handler

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

type TcpHandler struct {
}

func (t *TcpHandler) HandleRequest(conn net.Conn, sem chan bool) error {
	defer conn.Close()

	msgData := make([]byte, 2048)
	var len int
	var err error
	len, err = conn.Read(msgData)

	if err != nil {
		fmt.Println("error encountered while reading", err)
		conn.Close()
		<-sem
		return err
	}

	fmt.Printf("Data read from the client : %d\n", len)

	msg := string(msgData[:len])

	trimmedMsg := strings.TrimSpace(msg)
	trimmedMsg = strings.TrimSuffix(trimmedMsg, "\r\n")

	hexString := hex.EncodeToString([]byte(trimmedMsg))
	fmt.Printf("Command from Client: <<<<<<< %s >>>>>>\n", hexString)

	if trimmedMsg == "HELLO" {
		conn.Write([]byte(fmt.Sprintf("Hello from %s\n", conn.LocalAddr().String())))
		time.Sleep(3 * time.Second)
		<-sem
		return nil
	} else if trimmedMsg == "TIME" {
		conn.Write([]byte(fmt.Sprintf("Hello from %s\n", time.Now().String())))
	} else {
		conn.Write([]byte(fmt.Sprintf("invalid command : %s", trimmedMsg)))
	}
	<-sem
	return nil
}
