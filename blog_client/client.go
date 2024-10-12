package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")

	opts := grpc.WithTransportCredentials()

	cc, err := grpc.NewClient("localhost:5051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

}
