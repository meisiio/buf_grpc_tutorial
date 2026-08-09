package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/meisam/buf-tutorial/gen/go/proto/user/v1"
)

func main() {
	// 1. Set up a connection to the server.
	// Since we are running locally without TLS, we use insecure credentials.
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 2. Instantiate the generated client stub
	client := userv1.NewUserServiceClient(conn)

	// 3. Create a context with a 1-second timeout (Best practice in gRPC!)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 4. Make the RPC call for ID 101 (Success case)
	log.Println("--- Requesting User ID: 101 ---")
	req := &userv1.GetUserRequest{Id: 101}
	res, err := client.GetUser(ctx, req)
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("Success! Name: %s, Email: %s", res.GetUser().GetName(), res.GetUser().GetEmail())

	// 5. Make the RPC call for ID 999 (Error case)
	log.Println("--- Requesting User ID: 999 ---")
	req2 := &userv1.GetUserRequest{Id: 999}
	_, err2 := client.GetUser(ctx, req2)
	if err2 != nil {
		// Notice how gRPC translates the "NotFound" status code directly to an error here
		log.Printf("Expected Error Received: %v", err2)
	}
}
