package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	grpc_learn "github.com/Dubjay18/grpc-learn/blogpb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	grpc_learn.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	AuthorID string        `bson:"author_id"`
	Content  string        `bson:"content"`
	Title    string        `bson:"title"`
}

func (s *server) CreateBlog(ctx context.Context, req *grpc_learn.CreateBlogRequest) (*grpc_learn.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)

	}
	oid, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}
	return &grpc_learn.CreateBlogResponse{
		Blog: &grpc_learn.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		},
	}, nil
}

var collection *mongo.Collection

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	collection = client.Database("mydb").Collection("blog")
	fmt.Println("Blog Service Started")

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
	s.Stop()
	log.Println("Closing the listener")
	lis.Close()
	log.Println("End of program")
}
