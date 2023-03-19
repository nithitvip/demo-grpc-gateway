package service

import (
	pb "blog/proto/blog/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type BlogItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID  primitive.ObjectID `bson:"author_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func documentToBlog(data *BlogItem) *pb.Blog {
	return &pb.Blog{
		Id:        data.ID.Hex(),
		AuthorId:  data.AuthorID.Hex(),
		Title:     data.Title,
		Content:   data.Content,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}
}
