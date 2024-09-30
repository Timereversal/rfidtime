package main

import (
	"context"
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
var chanInventory = make(chan transport.RunnerData, 10)

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
	chipType := flag.String("chipType", "alienH3", "define chip type")

	flag.Parse()
	slog.Info(*addr)

	// establish Connection with Chafon decoder
	NewChafonConnection, err := transport.NewChafon(*addr, *chipType)
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

	b := sampling.Broker{StreamList: make(map[int32]chan transport.RunnerData)}

	// reading from one channel to deliver grpc message to server
	streamTagMax := make(chan transport.RunnerData)
	go func(streamMax <-chan transport.RunnerData) {
		for {
			select {
			case maxTag := <-streamMax:

				fmt.Printf("max tag id %+v \n ", maxTag)
				tagId := maxTag.TagID
				eventId := maxTag.EventId
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
	go func(in <-chan transport.RunnerData) {
		// how to read
		//for v := range in {
		//	fmt.Printf("%+v", v)
		//	epcS := fmt.Sprintf("%X", v.EPCData)
		//	slog.Info("log data structure", "epc", epcS, "rssi", v.RSSI)
		//}

		for {
			select {
			case runnerD := <-in:

				slog.Debug("%+v", runnerD)

				slog.Info("log data structure", "runnerID", runnerD.TagID, "eventID", runnerD.EventId, "rssi", runnerD.RSSI, "antenna", runnerD.Antenna)

				_, ok := b.StreamList[runnerD.TagID]
				if !ok {
					b.StreamGenerator(runnerD.TagID, streamTagMax)
				}
				// send TagInfo for Further procession [calculate  best sample ]
				go func() {
					b.StreamList[runnerD.TagID] <- runnerD
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
