package transport

type command struct {
	Len    byte
	Adr    byte //0x00 - 0xFE. 0xFF broadcast address
	Cmd    byte
	Data   []byte
	LsbCRC byte
	MsbCRC byte
}

type response struct {
	Len    byte
	Adr    byte
	ReCmd  byte
	Status byte
	Data   []byte
	LsbCRC byte
	MsbCRC byte
}
