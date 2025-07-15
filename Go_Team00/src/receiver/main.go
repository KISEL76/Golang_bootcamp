package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"Team00/frequencypb"
	receiver "Team00/receiver/internal"
)

func main() {
	k := flag.Float64("k", 3.0, "STD anomaly coefficient")
	flag.Parse()

	_ = godotenv.Load()
	dsn := os.Getenv("DSN")
	server_address := os.Getenv("SERVER_ADDRESS")
	db := receiver.NewDatabase(dsn)

	// gRPC connect
	conn, err := grpc.NewClient(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := frequencypb.NewFrequencyServiceClient(conn)

	stream, err := client.StreamFrequencies(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Error starting stream: %v", err)
	}

	detector := receiver.NewDetector(*k, 50)
	counter := 0

	log.Println("[START] Receiving frequency stream")

	for {
		entry, err := stream.Recv()
		if err != nil {
			log.Fatalf("Stream ended or error: %v", err)
		}

		value := entry.GetFrequency()
		detector.Process(value)
		counter++

		// периодический лог
		if counter%10 == 0 {
			fmt.Printf("[INFO] Processed %d values | Estimated mean=%.3f, std=%.3f\n",
				counter, detector.Mean(), detector.Std())
		}

		// проверка на аномалию
		if detector.IsReady() && detector.IsAnomaly(value) {
			fmt.Printf("[ANOMALY] session=%s | frequency=%.3f | mean=%.3f | std=%.3f | timestamp=%d\n",
				entry.GetSessionId(),
				value,
				detector.Mean(),
				detector.Std(),
				entry.GetTimestamp(),
			)

			db.SaveAnomaly(&receiver.Anomaly{
				SessionID: entry.GetSessionId(),
				Frequency: value,
				Timestamp: entry.GetTimestamp(),
			})
		}

		time.Sleep(10 * time.Millisecond)
	}
}
