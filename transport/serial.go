package transport

import (
	"go.bug.st/serial"
)

func stablishConnection(port string) (err error) {
	mode := &serial.Mode{
		BaudRate: 57600,
		Parity:   serial.NoParity,
		DataBits: 8,
	}

	return nil
}
