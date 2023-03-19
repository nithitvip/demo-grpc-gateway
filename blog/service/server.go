package service

import (
	pb "blog/proto/blog/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	BlogCollection *mongo.Collection

	pb.UnimplementedBlogServiceServer
	grpc_health_v1.UnimplementedHealthServer
}

func New(blogCollection *mongo.Collection) *Server {
	return &Server{
		BlogCollection: blogCollection,
	}
}
