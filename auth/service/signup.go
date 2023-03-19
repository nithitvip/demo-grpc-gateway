package service

import (
	pb "auth/proto/auth/v1"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Signup(ctx context.Context, request *pb.SignupRequest) (*pb.SignupResponse, error) {
	username := request.Username
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	password := request.Password
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	var result AuthItem
	err := s.AuthCollection.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	if err == nil {
		return nil, status.Error(codes.InvalidArgument, "username already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, status.Error(codes.Internal, err.Error())
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := AuthItem{
		Username: username,
		Password: hashedPassword,
	}

	res, err := s.AuthCollection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to oid: %v", err))
	}

	return &pb.SignupResponse{
		AccountId: oid.Hex(),
	}, nil
}
