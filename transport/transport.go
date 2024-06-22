package transport

import (
	"fmt"
	"io"
	"net"
)

func DialTcp(address string) error {
	//address := "192.168.1.200:27011"
	conn, err := net.Dial("tcp", address)
	// Data: 09000104fe00803280be
	bytesMessage := []byte{0x09, 0x00, 0x01, 0x04, 0xfe, 0x00, 0x80, 0x32, 0x80, 0xbe}
	_, err = conn.Write(bytesMessage)
	if err != nil {
		fmt.Printf("error during writting %s", err)
		return err
	}
	defer conn.Close()

	// 512KB (How to define a good buffer capacity?)
	// payload tcp size average ?
	buf := make([]byte, 1<<19)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Printf("Received: %d bytes %X \n", n, string(buf[:n]))
	}
	return nil
}

func deserialization(payload []byte, bytesNumber int) error {

	//var n int
	var totalBytes int

	for totalBytes < bytesNumber {
		lenPacket := int(payload[totalBytes])

		go handlePacket(payload[totalBytes : totalBytes+lenPacket])
		totalBytes += lenPacket + 1
	}

	return nil
}

func handlePacket(packet []byte) {
	fmt.Printf("%X", packet)
	return
}
