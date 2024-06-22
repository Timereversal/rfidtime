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

//Data: 15(len) [00](addr) [01](ReCMD)[03](Fruther data will be transfered)
//	    [01](antena)[01](only-1 tag )[0c]e28068940000500a9d2298c6  48257e
//		1500010301010c            e28068940000400a9d22a9eb  5c70d2
//		0d000103010104            111001ba465d27

//Data: 1500010301010c3034b708ac397e0000000649421ab7
//      0d000103010104111000a443ad25
//      15 00 01 03 01 01 0c 00 00 00 00 00 01 70 00 00 00 02 03 45 73 32

//Data: 150001030101 0c e28068940000400a9d22a9eb 5d f9c3  [22 bytes ] [0-21]
//      150001030101 0c e28068940000400a9d22acff 60 d3e2
//      150001030101 0c e28068940000500a9d2298c6 68 275f
//      150001030101 0c e28068940000400a9d2299e7 69 509b
//      150001030101 0c e28068940000500a9d22ad72 5a abaf
// 15 hex to dec 21
