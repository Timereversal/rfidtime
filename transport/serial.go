package transport

import (
	"fmt"
	"go.bug.st/serial"
)

func establishConnection(portA string) (*serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 57600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	serialPort, err := serial.Open(portA, mode)
	if err != nil {
		fmt.Printf("error during serial port creation  %s", err)
		return &serialPort, err
	}

	return &serialPort, nil
}
