package service

import (
	pb "auth/proto/auth/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	AuthCollection *mongo.Collection

	pb.UnimplementedAuthServiceServer
	grpc_health_v1.UnimplementedHealthServer
}

func New(authCollection *mongo.Collection) *Server {
	return &Server{
		AuthCollection: authCollection,
	}
}
