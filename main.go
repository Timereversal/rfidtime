package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"rfidtime/transport"
	"syscall"
)

// func calCRC16LSBMSB([]byte) (byte, byte) {
//
//		return 0, 0
//	}
var chanInventory = make(chan transport.DataInventory, 10)

func main() {
	//log.Print("Info message")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	file, err := os.Create("testing")
	if err != nil {
		slog.Error("Error opening file: ", err)
		panic(err)
	}

	//slog.Info("Info message")
	logger := slog.New(slog.NewTextHandler(file, nil))
	slog.SetDefault(logger)
	slog.Info("Hello world!")
	//slog.Info("hello", "count", 3)
	//var chafonCRC uint16 = 0x8408
	//chafon := crc16.Params{Poly: chafonCRC, Init: 0xFFFF}
	//
	addr := flag.String("address", "192.168.1.200:27011", "reader tcp/ip address:port")
	flag.Parse()
	slog.Info(*addr)

	NewChafonConnection, err := transport.NewChafon(*addr)
	if err != nil {
		slog.Info(err.Error())
	}
	go func() {
		<-c
		//cleanup()
		NewChafonConnection.SendCommand(transport.CmdModeAnswer)
		os.Exit(1)
	}()

	go func(in <-chan transport.DataInventory) {
		// how to read
		for v := range in {
			fmt.Printf("%+v", v)
			epcS := fmt.Sprintf("%X", v.EPCData)
			slog.Info("log data structure", "epc", epcS, "rssi", v.RSSI)
		}
	}(chanInventory)

	if err := NewChafonConnection.StartInventory(chanInventory); err != nil {
		slog.Info(err.Error())
	}

	defer NewChafonConnection.SendCommand(transport.CmdModeAnswer)

	//table := crc16.MakeTable(chafon)
	//crc := crc16.Checksum([]byte{0x15, 0x00, 0x01, 0x03, 0x01, 0x01, 0x0c, 0xe2, 0x80, 0x68, 0x94, 0x00, 0x00, 0x40, 0x0a, 0x9d, 0x22, 0xa9, 0xeb, 0x5d, 0xf9, 0xc3}, table)
	//fmt.Printf("CRC-16 MAXIM: %X\n", crc)
	//addr := "192.168.1.200:27011"
	//err := transport.DialTcp(addr)
	//if err != nil {
	//	fmt.Printf("%s", err)
	//}
	//conn, err := net.Dial("tcp", addr)
	//defer conn.Close()
	//if err != nil {
	//	fmt.Printf("error is %s", err)
	//}
	//fmt.Printf("starting connection")
	//buf := make([]byte, 1<<19) //512Kb
	//
	//for {
	//	n, err := conn.Read(buf)
	//	if err != nil {
	//		if err != io.EOF {
	//			fmt.Printf("%s", err)
	//		}
	//		break
	//	}
	//	fmt.Printf("read %d bytes, with data %s", n, buf[:n])
	//}
}
