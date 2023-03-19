package service

import (
	pb "auth/proto/auth/v1"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SignIn(ctx context.Context, request *pb.SignInRequest) (*pb.SignInResponse, error) {
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
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.InvalidArgument, "invalid user or password")
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if !CheckPasswordHash(password, result.Password) {
		return nil, status.Error(codes.InvalidArgument, "invalid user or password")
	}

	token, err := createToken(result.ID.Hex())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.SignInResponse{
		Token: token,
	}, nil
}
