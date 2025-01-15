package auth

import (
	"context"
	"fmt"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Auth) CreateUser(ctx context.Context, user CreateUserRequest) (CreateUserResponse, error) {
	const op = "gRPC Auth CreateUser"

	response, err := a.client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	})

	switch {
	case status.Code(err) == codes.AlreadyExists:
		return CreateUserResponse{}, UserAlreadyExistsErr
	case err != nil:
		return CreateUserResponse{}, fmt.Errorf(op, err)
	default:
		return CreateUserResponse{UserId: response.UserId}, nil
	}

}

func (a *Auth) Login(ctx context.Context, credentials LoginCredentials) (TokenResponse, error) {
	const op = "gRPC Auth Login"

	token, err := a.client.Login(ctx, &pb.LoginRequest{
		Username: credentials.Username,
		Password: credentials.Password,
	})

	switch {
	case status.Code(err) == codes.InvalidArgument:
		return TokenResponse{}, InvalidCredentialsErr
	case err != nil:
		return TokenResponse{}, fmt.Errorf(op, err)
	default:
		return TokenResponse{Value: token.Value}, nil
	}

}

func (a *Auth) VerifyToken(ctx context.Context, request TokenRequest) (VerifyTokenResponse, error) {
	const op = "gRPC Auth VerifyToken"

	verifyResponse, err := a.client.VerifyToken(ctx, &pb.TokenReuest{UserId: request.UserId, Token: request.Token})

	switch {
	case status.Code(err) == codes.InvalidArgument:
		return VerifyTokenResponse{}, InvalidCredentialsErr
	case err != nil:
		return VerifyTokenResponse{}, fmt.Errorf(op, err)
	default:
		return VerifyTokenResponse{Ok: verifyResponse.Ok}, nil
	}

}
