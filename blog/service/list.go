package service

import (
	pb "blog/proto/blog/v1"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *Server) ListBlog(ctx context.Context, in *emptypb.Empty) (*pb.ListBlogResponse, error) {
	log.Println("ListBlogRepeated was invoked")

	var blogs []*pb.Blog

	cur, err := s.BlogCollection.Find(context.Background(), primitive.D{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error: %v", err))
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &BlogItem{}
		err := cur.Decode(data)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("error: %v", err))
		}

		blogs = append(blogs, documentToBlog(data))
	}

	return &pb.ListBlogResponse{
		Blogs: blogs,
	}, nil
}
