package main

import (
	"context"
	"flag"
	"github.com/lmittmann/tint"
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

	t := time.Now()
	tFormat := t.Format("20060102150405")
	// Flag Attributes
	addr := flag.String("address", "192.168.1.200:27011", "reader tcp/ip address:port")
	chipType := flag.String("chipType", "alienH3", "define chip type")
	stage := flag.Int("stage", 1, "define stage to run")
	eventName := flag.String("eventName", "test_"+tFormat+".log", "define log filename")
	logType := flag.String("logType", "stdout", "define log type os.stdout or file")

	flag.Parse()

	switch *logType {
	case "stdout":
		//slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
		slog.SetDefault(slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.DateTime,
			}),
		))
	case "file":
		file, err := os.Create(*eventName + "_" + tFormat + ".log")
		if err != nil {
			slog.Error("Error opening file: ", err)
			panic(err)
		}
		defer file.Close()
		slog.SetDefault(slog.New(
			tint.NewHandler(file, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.DateTime,
			}),
		))

	}
	// define logger handler
	//logger := slog.New(slog.NewTextHandler(file, nil))
	//slog.SetDefault(logger)
	conngrpc, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error creating GRPCclient connection: ", err)

	}
	slog.Info("GRPC connection", "address", "127.0.0.1:8080", "status", "OK")
	defer conngrpc.Close()
	clientGRPC := reader.NewReaderClient(conngrpc)

	// establish Connection with Chafon decoder
	NewChafonConnection, err := transport.NewChafon(*addr, *chipType, *stage)
	if err != nil {
		slog.Error("New Reader Connection: ", "error", err.Error())
		panic(err)
	}
	stag := int32(NewChafonConnection.Stage)

	slog.Info("New Reader Connection: ", "address", *addr, "chipType", *chipType, "stage", stag)

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

	nSeconds := 30
	// reading from one channel to deliver grpc message to server
	streamTagMax := make(chan transport.RunnerData)
	go func(streamMax <-chan transport.RunnerData) {
		for {
			select {
			case maxTag := <-streamMax:

				slog.Debug("Pre-GRPC-send", "Stream Max Tag: ", maxTag)
				tagId := maxTag.TagID
				eventId := maxTag.EventId
				grpcTime := timestamppb.New(maxTag.Time)

				response, err := clientGRPC.Report(context.Background(), &reader.ReportRequest{TagId: tagId, EventId: eventId, RunnerTime: grpcTime, Stage: stag})
				if err != nil {

					slog.Error("GRPC-Send: ", "Error during grpc client report: ", err, "Runner-Data", maxTag)
				}
				slog.Info("GRPC-send", "Tag-ID", tagId, "Event-ID", eventId, "Runner-Time", maxTag.Time, "Stage", stag)
				slog.Debug("grpc Client", "response", response)
			// No receive any Runner Data Best Sample During n seconds,
			case <-time.After(time.Duration(nSeconds) * time.Second):
				slog.Info("Runner Data Best Sample not recevied during ", "seconds", nSeconds)

			}
		}
	}(streamTagMax)
	// Listening Channel info to generate log for each bib-tag, RunnerData struct send by handlepacket in transport package
	// assign TagInfo a corresponding Channel for specific procession (max RSSI value)
	go func(in <-chan transport.RunnerData) {

		for {
			select {
			case runnerD := <-in:

				slog.Debug("Runner-Data", "%+v", runnerD)

				slog.Info("Runner-Data", "runnerID", runnerD.TagID, "eventID", runnerD.EventId, "rssi", runnerD.RSSI, "antenna", runnerD.Antenna)

				// if there is not a channel for best sample analisis in StreamList, it will generate a new on.
				_, ok := b.StreamList[runnerD.TagID]
				if !ok {
					b.StreamGenerator(runnerD.TagID, streamTagMax)
					slog.Debug("Runner-Tag-Id Analysis: ", "Create ID Pool", runnerD.TagID)
				}
				// send TagInfo for Further procession [calculate  best sample ]
				go func() {
					slog.Debug("Runner-Tag-Id Analysis: ", "Send Runner Data to Pool", runnerD)
					b.StreamList[runnerD.TagID] <- runnerD
				}()

			}
		}
	}(chanInventory)

	// Start Inventory.
	if err := NewChafonConnection.StartInventory(chanInventory); err != nil {
		slog.Error("Start-Service: ", "error", err.Error())
	}
	b.Wg.Wait()

}
