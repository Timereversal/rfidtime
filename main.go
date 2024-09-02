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

// Channel to send logs for each goroutine[bib-tag] to a log process
var chanInventory = make(chan transport.DataInventory, 10)

func main() {

	// define log file
	file, err := os.Create("testing")
	if err != nil {
		slog.Error("Error opening file: ", err)
		panic(err)
	}

	// define logger handler
	logger := slog.New(slog.NewTextHandler(file, nil))
	slog.SetDefault(logger)

	// Chafon address decoder
	addr := flag.String("address", "192.168.1.200:27011", "reader tcp/ip address:port")
	flag.Parse()
	slog.Info(*addr)

	// establish Connection with Chafon decoder
	NewChafonConnection, err := transport.NewChafon(*addr)
	if err != nil {
		slog.Info(err.Error())
	}

	// listening for Ctr+C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		//cleanup()
		_ = NewChafonConnection.SendCommand(transport.CmdModeAnswer)
		os.Exit(1)
	}()

	// Listening Channel info to generate log for each bib-tag
	go func(in <-chan transport.DataInventory) {
		// how to read
		for v := range in {
			fmt.Printf("%+v", v)
			epcS := fmt.Sprintf("%X", v.EPCData)
			slog.Info("log data structure", "epc", epcS, "rssi", v.RSSI)
		}
	}(chanInventory)

	// Start Inventory.
	if err := NewChafonConnection.StartInventory(chanInventory); err != nil {
		slog.Info(err.Error())
	}

}
