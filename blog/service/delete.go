package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "blog/proto/blog/v1"
)

func (s *Server) DeleteBlog(ctx context.Context, in *pb.DeleteBlogRequest) (*emptypb.Empty, error) {

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

	_, err = s.BlogCollection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot Delete: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) validateBlogOwner(ctx context.Context, blogId primitive.ObjectID, accountId string) error {

	data := &BlogItem{}

	filter := bson.M{"_id": blogId}

	res := s.BlogCollection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return status.Errorf(codes.NotFound, "Cannot find blog with ID provided")
	}

	if data.AuthorID.Hex() != accountId {
		return status.Errorf(codes.PermissionDenied, "not owner of this blog")
	}

	return nil
}
