package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"connectrpc.com/connect"
	userv1 "github.com/meisam/buf-tutorial/gen/go/proto/user/v1"
	userv1connect "github.com/meisam/buf-tutorial/gen/go/proto/user/v1/v1connect"
)

type connectServer struct{}

func (s *connectServer) GetUser(ctx context.Context, req *connect.Request[userv1.GetUserRequest]) (*connect.Response[userv1.GetUserResponse], error) {
	log.Printf("Received Connect request for user ID: %d", req.Msg.GetId())
	if req.Msg.GetId() != 101 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user with this id not found"))
	}

	res := connect.NewResponse(&userv1.GetUserResponse{
		User: &userv1.User{
			Id:    101,
			Name:  "Meisam",
			Email: "[EMAIL_ADDRESS]",
		},
	})
	return res, nil
}
func (s *connectServer) ListUsers(
	ctx context.Context,
	req *connect.Request[userv1.ListUserRequest],
	stream *connect.ServerStream[userv1.User],
) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("stream not built yet"))
}
func main() {
	path, handler := userv1connect.NewUserServiceHandler(&connectServer{})
	http.Handle(path, handler)

	log.Printf("ConnectRPC Server listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
