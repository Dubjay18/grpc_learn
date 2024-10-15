package main

import (
	"context"
	"fmt"
	"log"

	grpc_learn "github.com/Dubjay18/grpc-learn/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")

	opts := grpc.WithTransportCredentials(nil)

	cc, err := grpc.NewClient("localhost:5051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := grpc_learn.NewBlogServiceClient(cc)
	fmt.Println("Creating blog")
	blog := &grpc_learn.Blog{
		AuthorId: "Jay",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &grpc_learn.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error %v", err)
	}

	fmt.Printf("Blog created successfully %v", createBlogRes)

}
