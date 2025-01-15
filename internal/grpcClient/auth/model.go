package auth

import (
	"context"
	"fmt"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
	"google.golang.org/grpc"
)

type Authorizer interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (CreateUserResponse, error)
	Login(ctx context.Context, credentials LoginCredentials) (TokenResponse, error)
	VerifyToken(ctx context.Context, token TokenRequest) (VerifyTokenResponse, error)
}

var UserAlreadyExistsErr = fmt.Errorf("already exists")
var InvalidCredentialsErr = fmt.Errorf("invalid credentials")

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

type TokenResponse struct {
	Value string
}

type TokenRequest struct {
	UserId string
	Token  string
}

type VerifyTokenResponse struct {
	Ok bool
}
