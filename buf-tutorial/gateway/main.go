package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/meisam/buf-tutorial/gen/go/proto/user/v1"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. Create a new gRPC-Gateway ServeMux (this is the REST router!)
	mux := runtime.NewServeMux()

	// 2. Set up the connection to our actual backend gRPC Server
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 3. Register our generated Gateway proxy code
	err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// 4. Start the HTTP server!
	log.Println("REST Gateway listening on port 8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve gateway: %v", err)
	}
}
