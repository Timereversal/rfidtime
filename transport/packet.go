package transport

import (
	"fmt"
	"strconv"
	"time"
)

type ParseErr struct {
	Message string
	Err     error
}

func (e ParseErr) Error() string {
	return e.Message
}

func (e ParseErr) Unwrap() error {
	return e.Err
}

type Command struct {
	Len    byte
	Adr    byte //0x00 - 0xFE. 0xFF broadcast address
	Cmd    byte
	Data   []byte
	LsbCRC byte
	MsbCRC byte
}

type Response struct {
	Len    byte
	Adr    byte
	ReCmd  byte
	Status byte
	Data   []byte
	LsbCRC byte
	MsbCRC byte
}

type TagInfo struct {
	Ant     byte
	Num     byte
	EPCLen  byte
	EPCData []byte
	RSSI    int
	Phase   []byte
	Freq    []byte
	LSB     byte
	MSB     byte
	Time    time.Time
}

// EPCData encode User-Event Information
type EPCData struct {
	RunnerID [2]byte // 2 bytes enough,
	EventID  [2]byte //
	OrgID    [2]byte
}

// RunnerData store Runner information
// TagID-EventId is unique
type RunnerData struct {
	TagID   int32
	EventId int32
	Stage   int
	RSSI    int
	Antenna int
	Time    time.Time
}

//type EPC struct {
//	EPCLen byte
//	EPCData    []byte
//	RSSI    byte
//	Phase   []byte
//	Freq    []byte
//}

func ParseResponse(r Response, typ string) (RunnerData, error) {
	// Check if info is valid
	//runnerData := RunnerData{}
	antenna := int(r.Data[0])
	epcData := r.Data[2 : 2+int(r.Data[1])]
	rssi := int(r.Data[2+int(r.Data[1])])
	switch typ {
	case "alienH3":
		//alienH3 chip data len - 12 bytes [24 characters].
		// last 6 characters tagID
		if len(epcData) != 12 {
			return RunnerData{}, ParseErr{
				Message: "invalid length of alienH3",
			}
		}
		temp := fmt.Sprintf("%X", epcData)
		tagId, err := strconv.Atoi(temp[len(temp)-6:])

		if err != nil {
			return RunnerData{}, ParseErr{
				Message: fmt.Sprintf("invalid tagId %s", temp[len(temp)-6:]),
				Err:     err,
			}

		}
		eventId, err := strconv.Atoi(temp[len(temp)-12 : len(temp)-6])
		if err != nil {
			return RunnerData{}, ParseErr{
				Message: fmt.Sprintf("invalid eventId %s", temp[len(temp)-12:len(temp)-6]),
				Err:     err,
			}
		}

		return RunnerData{TagID: int32(tagId), EventId: int32(eventId), RSSI: rssi, Antenna: antenna}, nil

	}

	return RunnerData{}, fmt.Errorf("error during parsing %+v  chipType %s", r, typ)
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
