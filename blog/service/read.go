package service

import (
	pb "blog/proto/blog/v1"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func (s *Server) ReadBlog(ctx context.Context, in *pb.ReadBlogRequest) (*pb.ReadBlogResponse, error) {
	log.Printf("ReadBlog was invoked with %v\n", in)
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot Parse ID")
	}

	data := &BlogItem{}
	filter := bson.M{"_id": oid}

	res := s.BlogCollection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "Cannot find blog with ID provided")
		} else {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	return &pb.ReadBlogResponse{
		Id:        data.ID.Hex(),
		AuthorId:  data.AuthorID.Hex(),
		Title:     data.Title,
		Content:   data.Content,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}, nil
}
