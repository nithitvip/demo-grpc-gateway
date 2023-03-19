package service

import (
	pb "blog/proto/blog/v1"
	"context"
	"fmt"
	fieldmaskutils "github.com/mennanov/fieldmask-utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strings"
	"time"
)

func (s *Server) PartialUpdateBlog(ctx context.Context, in *pb.PartialUpdateBlogRequest) (*emptypb.Empty, error) {
	log.Printf("PartialUpdateBlog was invoked with %v\n", in)
	log.Printf("field mask %s", in.GetFieldMask())
	accountId, ok := getAccountId(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("account_id metadata is required"))
	}

	oid, err := primitive.ObjectIDFromHex(in.UpdateReq.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot Parse ID")
	}

	// check if submitted id is owner of this blog
	err = s.validateBlogOwner(ctx, oid, accountId)
	if err != nil {
		return nil, err
	}

	updateFields := make(map[string]interface{})
	mask, _ := fieldmaskutils.MaskFromProtoFieldMask(in.FieldMask, naming)

	err = fieldmaskutils.StructToMap(mask, in.UpdateReq, updateFields)

	updateFields["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{}
	for key, value := range updateFields {
		update[strings.ToLower(key)] = value
	}
	res, err := s.BlogCollection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: %v", err)
	}

	if res.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Cannot find blog with Id")
	}
	return &emptypb.Empty{}, nil
}

func naming(s string) string {
	return strings.Title(s)
}
