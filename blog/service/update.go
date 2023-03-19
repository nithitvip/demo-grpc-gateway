package service

import (
	pb "blog/proto/blog/v1"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *Server) UpdateBlog(ctx context.Context, in *pb.UpdateBlogRequest) (*emptypb.Empty, error) {
	log.Printf("UpdateBlog was invoked with %v\n", in)

	accountId, ok := getAccountId(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("account_id metadata is required"))
	}

	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot Parse ID")
	}

	// check if submitted id is owner of this blog
	err = s.validateBlogOwner(ctx, oid, accountId)
	if err != nil {
		return nil, err
	}

	acct, err := primitive.ObjectIDFromHex(accountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot Parse ID")
	}

	data := &BlogItem{AuthorID: acct, Title: in.Title, Content: in.Content}

	res, err := s.BlogCollection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": data},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	if res.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Cannot find blog with Id")
	}
	return &emptypb.Empty{}, nil
}
