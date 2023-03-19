package main

import (
	"context"
	"fmt"
	"gateway/interceptor"
	authpb "gateway/proto/auth/v1"
	blogpb "gateway/proto/blog/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

func main() {
	gwmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(CustomMatcher))
	err := registerAuthService(gwmux)
	if err != nil {
		log.Fatal(err)
	}

	err = registerBlogService(gwmux)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	mux.Handle("/openapiv2/", http.StripPrefix("/openapiv2", http.FileServer(http.Dir("openapiv2"))))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", http.FileServer(http.Dir("dist"))))

	gatewayAddr := "0.0.0.0:8080"
	gwServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: mux,
	}
	log.Println("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	log.Fatalln(gwServer.ListenAndServe())
}

func registerAuthService(gwmux *runtime.ServeMux) error {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.LogUnaryClientInterceptor()),
	}
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(context.Background(), gwmux, "localhost:50051", dialOpts)
	if err != nil {
		return fmt.Errorf("failed to register auth: %w", err)
	}
	return nil
}

func registerBlogService(gwmux *runtime.ServeMux) error {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.LogUnaryClientInterceptor()),
		grpc.WithUnaryInterceptor(interceptor.AuthUnaryClientInterceptor()),
	}
	err := blogpb.RegisterBlogServiceHandlerFromEndpoint(context.Background(), gwmux, "localhost:50052", dialOpts)
	if err != nil {
		return fmt.Errorf("failed to register blog: %w", err)
	}
	return nil
}

func CustomMatcher(key string) (string, bool) {
	return key, false
}
