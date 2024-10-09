package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	grpc_learn "github.com/Dubjay18/grpc-learn/blogpb"
	"google.golang.org/grpc"
)

type server struct {
	grpc_learn.UnimplementedBlogServiceServer
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	grpc_learn.RegisterBlogServiceServer(s, &server{})

	go func() {
		log.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	log.Println("Stopping the server")

}
