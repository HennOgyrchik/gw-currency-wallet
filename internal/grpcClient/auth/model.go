package auth

import (
	"context"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
	"google.golang.org/grpc"
)

type Authorizer interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (CreateUserResponse, error)
	Login(ctx context.Context, credentials LoginCredentials) (Token, error)
	VerifyToken(ctx context.Context, token Token) (VerifyTokenResponse, error)
}

type Auth struct {
	url    string
	conn   *grpc.ClientConn
	client pb.AuthorizationClient
}

type CreateUserRequest struct {
	Username string
	Password string
	Email    string
}

type CreateUserResponse struct {
	UserId string
}

type LoginCredentials struct {
	Username string
	Password string
}

type Token struct {
	Value string
}

type VerifyTokenResponse struct {
	Ok bool
}
