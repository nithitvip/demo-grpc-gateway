package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strings"

	"gateway/jwt"
)

func AuthUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("method:", method)
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
		}

		authHeader, ok := md["authorization"]
		if !ok {
			return status.Errorf(codes.Unauthenticated, "Token is not provided")
		}

		if len(authHeader) == 0 {
			return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
		}

		accountId, err := jwt.VerifyAndGetId(strings.TrimPrefix(authHeader[0], "Bearer "))
		if err != nil {
			return err
		}

		log.Printf("method: %s with account_id: %s", method, accountId)
		ctx = metadata.AppendToOutgoingContext(ctx, "account_id", accountId)

		err = invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
