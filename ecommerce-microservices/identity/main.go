package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"connectrpc.com/connect"
	identityv1 "github.com/meisam/ecommerce/gen/go/proto/identity/v1"
	"github.com/meisam/ecommerce/gen/go/proto/identity/v1/v1connect"
)

type identityServer struct{}

func (S *identityServer) CreateUser(
	ctx context.Context,
	req *connect.Request[identityv1.CreateUserRequest],
) (*connect.Response[identityv1.CreateUserResponse], error) {

	log.Printf("creating user : %s", req.Msg.GetEmail())

	res := connect.NewResponse(&identityv1.CreateUserResponse{
		UserId: 99,
	})

	return res, nil

}

func (S *identityServer) Login(
	ctx context.Context,
	req *connect.Request[identityv1.LoginRequest],
) (*connect.Response[identityv1.LoginResponse], error) {
	email, password := req.Msg.GetEmail(), req.Msg.GetPassword()

	if email == "admin@example.com" && password == "secret" {
		jwt_token := "fake_jwt_token"
		return connect.NewResponse(&identityv1.LoginResponse{
			JwtToken: jwt_token,
		}), nil

	}
	return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))

}

func main() {
	path, handler := v1connect.NewIdentityServiceHandler(&identityServer{})
	http.Handle(path, handler)
	log.Printf("server is running on 8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
