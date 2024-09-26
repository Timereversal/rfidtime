package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"os"
	"os/signal"
	"rfidtime/reader"
	"rfidtime/sampling"
	"rfidtime/transport"
	"syscall"
	"time"
)

// Channel to send logs for each goroutine[bib-tag] to a log process
var chanInventory = make(chan transport.TagInfo, 10)

func main() {

	//var opts []grpc.DialOption
	//conn, err := grpc.NewClient("127.0.0.1:8080", opts...)
	//client := reader.NewReaderClient(conn)

	// define log file
	file, err := os.Create("testing")
	if err != nil {
		slog.Error("Error opening file: ", err)
		panic(err)
	}
	conngrpc, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error creating GRPCclient connection: ", err)
	}
	defer conngrpc.Close()
	clientGRPC := reader.NewReaderClient(conngrpc)

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

	b := sampling.Broker{StreamList: make(map[string]chan transport.TagInfo)}

	// reading from one channel to deliver grpc message to server
	streamTagMax := make(chan transport.TagInfo)
	go func(streamMax <-chan transport.TagInfo) {
		for {
			select {
			case maxTag := <-streamMax:

				fmt.Printf("max tag id %X \n ", maxTag.EPCData)
				tagId := int32(binary.BigEndian.Uint32(maxTag.EPCData[len(maxTag.EPCData)-4:]))
				eventId := int32(binary.BigEndian.Uint32(maxTag.EPCData[len(maxTag.EPCData)-7:]))
				grpcTime := timestamppb.New(maxTag.Time)

				response, err := clientGRPC.Report(context.Background(), &reader.ReportRequest{TagId: tagId, EventId: eventId, RunnerTime: grpcTime})
				if err != nil {
					fmt.Println("Error during grpc client report: ", err)
				}
				fmt.Println(response)
			case <-time.After(30 * time.Second):
				fmt.Println("timexxxxxx")

			}
		}
	}(streamTagMax)
	// Listening Channel info to generate log for each bib-tag
	// assign TagInfo a corresponding Channel for specific procession (max RSSI value)
	go func(in <-chan transport.TagInfo) {
		// how to read
		//for v := range in {
		//	fmt.Printf("%+v", v)
		//	epcS := fmt.Sprintf("%X", v.EPCData)
		//	slog.Info("log data structure", "epc", epcS, "rssi", v.RSSI)
		//}

		for {
			select {
			case tagInfo := <-in:
				slog.Debug("%+v", tagInfo)
				epcS := fmt.Sprintf("%X", tagInfo.EPCData)

				n := len(tagInfo.EPCData)
				tagInfoID := fmt.Sprintf("%X", tagInfo.EPCData[n-4:])
				slog.Info("log data structure", "epc", epcS, "rssi", tagInfo.RSSI, "tagID", tagInfoID)

				_, ok := b.StreamList[tagInfoID]
				if !ok {
					b.StreamGenerator(tagInfoID, streamTagMax)
				}
				// send TagInfo for Further procession [calculate  best sample ]
				go func() {
					b.StreamList[tagInfoID] <- tagInfo
				}()

			}
		}
	}(chanInventory)

	// Start Inventory.
	if err := NewChafonConnection.StartInventory(chanInventory); err != nil {
		slog.Info(err.Error())
	}
	b.Wg.Wait()

}
