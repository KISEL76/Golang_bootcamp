package main

import (
	"Team00/frequencypb"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	frequencypb.UnimplementedFrequencyServiceServer
	rng *rand.Rand
}

func (s *server) StreamFrequencies(_ *emptypb.Empty, stream frequencypb.FrequencyService_StreamFrequenciesServer) error {
	sessionID := uuid.New().String()
	mean := s.rng.Float64()*20 - 10
	stddev := s.rng.Float64()*1.2 + 0.3
	log.Printf("New stream: session=%s mean=%.2f stddev=%.2f", sessionID, mean, stddev)

	for i := 0; i < 200; i++ {
		frequency := s.rng.NormFloat64()*stddev + mean
		timestamp := time.Now().Unix()

		entry := &frequencypb.FrequencyEntry{
			SessionId: sessionID,
			Frequency: frequency,
			Timestamp: timestamp,
		}

		if err := stream.Send(entry); err != nil {
			log.Printf("Failed to send: %v", err)
			return err
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	frequencypb.RegisterFrequencyServiceServer(s, &server{
		rng: rng,
	})

	log.Println("gRPC server listening on :50051")
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
