package service

import (
	pb "blog/proto/blog/v1"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func (s *Server) CreateBlog(ctx context.Context, request *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	log.Printf("CreateBlog was invoked with %v\n", request)

	accountId, ok := getAccountId(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("account_id metadata is required"))
	}

	acctId, err := primitive.ObjectIDFromHex(accountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot Parse ID")
	}
	data := BlogItem{
		AuthorID:  acctId,
		Title:     request.Title,
		Content:   request.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := s.BlogCollection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to oid: %v", err))
	}

	return &pb.CreateBlogResponse{
		Id: oid.Hex(),
	}, nil
}

func getAccountId(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	accountId := md.Get("account_id")
	if len(accountId) == 0 {
		return "", false
	}
	return accountId[0], true
}
