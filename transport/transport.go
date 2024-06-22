package transport

import (
	"fmt"
	"io"
	"net"
	"time"
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
		err = deserialization(buf[:n], n)
		if err != nil {
			return fmt.Errorf("error deserializing %s", err)
		}
	}
	return nil
}

func deserialization(payload []byte, bytesNumber int) error {

	//var n int
	var totalBytes int

	for totalBytes < bytesNumber {
		lenPacket := int(payload[totalBytes])

		go handlePacket(payload[totalBytes : totalBytes+lenPacket+1])
		//go handlePacket(payload[totalBytes : totalBytes+lenPacket])
		totalBytes += lenPacket + 1
	}

	return nil
}

func handlePacket(packet []byte) {
	//fmt.Printf("%X ", packet)
	dt := time.Now()
	// packet analysis response.
	//MSB LSB are not included
	packetR := Response{Len: packet[0], Adr: packet[1], ReCmd: packet[2], Status: packet[3], Data: packet[4:len(packet)]}

	// handling Response Command of 0x01 [Tag Inventory request]
	if packetR.ReCmd == 0x01 {
		switch packetR.Status {
		case 0x01:
			fmt.Println("tag inventory command succesfull delivered, reader will transmit")

		case 0x02:
			fmt.Println("tag inventory command, reader fails to complete the inventory within the predefined inventory time.")
		case 0x03:
			// EPCData is Ok if bit 6 and bit7 of PacketR.Data[2] are 0
			epcInfo := DataInventory{Ant: packetR.Data[0],
				EPCData: packetR.Data[3 : 3+int(packetR.Data[2])],
				RSSI:    packetR.Data[3+int(packetR.Data[2])]}
			fmt.Printf("%d,epc:%X,rssi:%d,%X Packet being transmitted at %s \n", packetR.Len, epcInfo.EPCData, epcInfo.RSSI, packetR.Data, dt.Format("01-02-2006 15:04:05"))
		case 0xF8:
			fmt.Println("Antenna Error Detected")
		}
	}
	return
}
