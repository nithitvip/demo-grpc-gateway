package main

import (
	"context"
	"fmt"
	"gateway/interceptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strings"

	authpb "gateway/proto/auth/v1"
	blogpb "gateway/proto/blog/v1"
)

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
//func getOpenAPIHandler() http.Handler {
//	mime.AddExtensionType(".svg", "image/svg+xml")
//	// Use subdirectory in embedded files
//	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
//	if err != nil {
//		panic("couldn't create sub filesystem: " + err.Error())
//	}
//	return http.FileServer(http.FS(subFS))
//}

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

	//oa := getOpenAPIHandler()
	gatewayAddr := "0.0.0.0:8080"
	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				gwmux.ServeHTTP(w, r)
				return
			}
			//		oa.ServeHTTP(w, r)
		}),
	}
	log.Println("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	log.Fatalln(gwServer.ListenAndServe())
}

func registerAuthService(gwmux *runtime.ServeMux) error {
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(context.Background(), gwmux, "localhost:50051", dialOpts)
	if err != nil {
		return fmt.Errorf("failed to register auth: %w", err)
	}
	return nil
}

func registerBlogService(gwmux *runtime.ServeMux) error {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
