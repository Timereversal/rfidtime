package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	addr := "192.168.1.200:27011"
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()
	if err != nil {
		fmt.Printf("error is %s", err)
	}
	fmt.Printf("starting connection")
	buf := make([]byte, 1<<19) //512Kb

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("%s", err)
			}
			break
		}
		fmt.Printf("read %d bytes, with data %s", n, buf[:n])
	}
}
