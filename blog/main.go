package main

import (
	pb "blog/proto/blog/v1"
	"blog/service"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

var addr = "0.0.0.0:50052"

func main() {
	client, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	collection := client.Database("blogdb").Collection("blog")

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	defer lis.Close()

	serv := service.New(collection)

	log.Printf("Listening at %s\n", addr)
	s := grpc.NewServer()

	grpc_health_v1.RegisterHealthServer(s, serv)
	pb.RegisterBlogServiceServer(s, serv)

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

func connectDB() (*mongo.Client, error) {
	db := getEnv("MONGO_URI", "localhost:27017")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@" + db + "/"))
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
