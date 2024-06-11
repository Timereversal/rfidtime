package transport

import (
	"fmt"
	"io"
	"net"
)

func DialTcp() {
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
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
		fmt.Printf("Received: %d bytes %s", n, string(buf[:n]))
	}

}
