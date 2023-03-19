package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

func LogUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("request to method:", method)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
