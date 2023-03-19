package main

//
//import (
//	"context"
//	"log"
//	"net/http"
//
//	authpb "gateway/proto/auth/v1"
//
//	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//)
//
//func main() {
//	option := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
//
//	gwmux := runtime.NewServeMux()
//
//	err := authpb.RegisterAuthServiceHandlerFromEndpoint(context.Background(), gwmux, "localhost:50051", option)
//	if err != nil {
//		log.Fatalln("Failed to register gateway:", err)
//	}
//
//	gwServer := &http.Server{
//		Addr:    ":8080",
//		Handler: gwmux,
//	}
//
//	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8080")
//	log.Fatalln(gwServer.ListenAndServe())
//}
