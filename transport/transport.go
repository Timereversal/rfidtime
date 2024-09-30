package transport

import (
	"fmt"
	"io"
	"net"
	"time"
)

var cmdAnswerMode = []byte{0x09, 0x00, 0x01, 0x04, 0xfe, 0x00, 0x80, 0x32, 0x80, 0xbe}

var cmdRealTimeInventoryParam = []byte{}
var cmdModeRealTimeInventory = []byte{0x05, 0x00, 0x76, 0x01, 0xeb, 0xd8}
var CmdModeAnswer = []byte{0x05, 0x00, 0x76, 0x00, 0x62, 0xc9}

// channel to communicate with logger goroutine (struct TagInfo transmitted)

type ChafonInterface interface {
	SendCommand([]byte) error
	startInventory() error
}

type Chafon struct {
	connection *net.Conn
	ChipType   string
}

// NewChafon generate a Chafon struct pointer where address is chafon reader ip
// address format example 192.168.1.200:27011
func NewChafon(address string, chipType string) (*Chafon, error) {
	conn, err := net.Dial("tcp", address)
	//defer conn.Close()
	if err != nil {
		return &Chafon{}, nil
	}

	cf := Chafon{
		connection: &conn,
		ChipType:   chipType,
	}
	return &cf, nil
}

// SendCommand delivery commands trough Chafon net connection
func (cf *Chafon) SendCommand(cmd []byte) error {
	// Which kind of errors are considered in Net Conn Write ?
	_, err := (*cf.connection).Write(cmd)
	//_, err = cf.connection.Write(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (cf *Chafon) ReceiveCommand() ([]byte, error) {
	return []byte{}, nil
}

// StartInventory - there are several ways to perform Inventory
// current mode only support RealTimeInventory
func (cf *Chafon) StartInventory(out chan<- RunnerData) error {
	// byte qty per rfid AlienH3  22bytes
	err := cf.SendCommand(cmdModeRealTimeInventory)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	for {
		n, err := (*cf.connection).Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		fmt.Printf("Received: %d bytes %X \n", n, string(buf[:n]))
		//err = deserialization(buf[:n], n)
		err = cf.deserialization(buf[:n], n, out)
		if err != nil {
			return fmt.Errorf("error deserializing %s", err)
		}

	}
	return nil
}

// deserialization, TCP DATA contains more than 1 EPC information. this function generate a new goroutine to process
// each EPC information using the handlePacket function
func (cf *Chafon) deserialization(payload []byte, bytesNumber int, out chan<- RunnerData) error {
	var totalBytes int

	for totalBytes < bytesNumber {
		lenPacket := int(payload[totalBytes])

		//go handlePacket(payload[totalBytes : totalBytes+lenPacket+1])

		go cf.handlePacket(payload[totalBytes:totalBytes+lenPacket+1], out)
		totalBytes += lenPacket + 1
	}

	return nil
}

func (cf *Chafon) handlePacket(packet []byte, out chan<- RunnerData) {

	dt := time.Now()
	// packet analysis response.
	//MSB LSB are not included
	// Packet parsing
	packetR := Response{Len: packet[0], Adr: packet[1], ReCmd: packet[2], Status: packet[3], Data: packet[4:len(packet)]}

	// handling response of RealTime inventory [Handle set mode Realtime Inventory]
	if packetR.ReCmd == 0xee { // set mode RealTime Inventory
		switch packetR.Status {
		case 0x00:
			// Detected appropiate tag
			//fmt.Println("Detected Appropiate tag")
			epcInfo := TagInfo{
				Ant:     packetR.Data[0],
				EPCData: packetR.Data[2 : 2+int(packetR.Data[1])],
				RSSI:    int(packetR.Data[2+int(packetR.Data[1])]),
				Time:    dt,
			}
			runnerD, err := ParseResponse(packetR, cf.ChipType)
			if err != nil {
				fmt.Printf("Error parsing response: %s\n", err)
				return
			}
			runnerD.Time = dt
			fmt.Printf("%d,epc:%X,rssi:%d,%X Packet being transmitted at %s \n", packetR.Len, epcInfo.EPCData, epcInfo.RSSI, packetR.Data, dt.Format("01-02-2006 15:04:05"))
			fmt.Printf("Runner-ID: %d, EventID: %d, RSSI: %d, time: %s \n", runnerD.TagID, runnerD.EventId, runnerD.RSSI, dt.Format("01-02-2006 15:04:05"))
			out <- runnerD
		}
		return
	}
	// handling Response Command of 0x01 [Answer Mode - Tag Inventory request]
	if packetR.ReCmd == 0x01 {
		switch packetR.Status {
		case 0x01:
			fmt.Println("tag inventory command succesfull delivered, reader will transmit")
			// Send request Inventory to start again .
			err := cf.SendCommand(cmdAnswerMode)
			if err != nil {
				fmt.Println("error during HandlePacket Case 0x01")
			}
			// Send continuosly Inventory commands or stop gracefully
			//err := ChafonInterface.SendCommand(cmdAnswerMode)
			//if err != nil {
			//	fmt.Printf("error during sending command to Chafon %s", err)
			//}
		case 0x02:
			fmt.Println("tag inventory command, reader fails to complete the inventory within the predefined inventory time.")
		case 0x03:
			// EPCData is Ok if bit 6 and bit7 of PacketR.Data[2] are 0
			epcInfo := TagInfo{Ant: packetR.Data[0],
				EPCData: packetR.Data[3 : 3+int(packetR.Data[2])],
				RSSI:    int(packetR.Data[3+int(packetR.Data[2])])}
			fmt.Printf("%d,epc:%X,rssi:%d,%X Packet being transmitted at %s \n", packetR.Len, epcInfo.EPCData, epcInfo.RSSI, packetR.Data, dt.Format("01-02-2006 15:04:05"))
		case 0xF8:
			fmt.Println("Antenna Error Detected")
		}
		return
	}

}

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
		//slog.Debug("Received: %d bytes %X \n", n, string(buf[:n]))
		err = deserialization(buf[:n], n)
		if err != nil {
			return fmt.Errorf("error deserializing %s", err)
		}
	}
	return nil
}

// func deserialization split stream of serial packets.
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
			// Send continuosly Inventory commands or stop gracefully
			//err := ChafonInterface.SendCommand(cmdAnswerMode)
			//if err != nil {
			//	fmt.Printf("error during sending command to Chafon %s", err)
			//}
		case 0x02:
			fmt.Println("tag inventory command, reader fails to complete the inventory within the predefined inventory time.")
		case 0x03:
			// EPCData is Ok if bit 6 and bit7 of PacketR.Data[2] are 0
			epcInfo := TagInfo{Ant: packetR.Data[0],
				EPCData: packetR.Data[3 : 3+int(packetR.Data[2])],
				RSSI:    int(packetR.Data[3+int(packetR.Data[2])]),
				Time:    dt,
			}
			fmt.Printf("%d,epc:%X,rssi:%d,%X Packet being transmitttttttted at %s \n", packetR.Len, epcInfo.EPCData, epcInfo.RSSI, packetR.Data, dt.Format("01-02-2006 15:04:05"))
		case 0xF8:
			fmt.Println("Antenna Error Detected")
		}
	}
	return
}
