package auth

import (
	"context"
	"fmt"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
)

func (a *Auth) CreateUser(ctx context.Context, user CreateUserRequest) (CreateUserResponse, error) {
	const op = "gRPC Auth CreateUser"

	response, err := a.client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	})

	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return CreateUserResponse{UserId: response.UserId}, err
}

func (a *Auth) Login(ctx context.Context, credentials LoginCredentials) (Token, error) {
	const op = "gRPC Auth Login"

	token, err := a.client.Login(ctx, &pb.LoginRequest{
		Username: credentials.Username,
		Password: credentials.Password,
	})
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return Token{Value: token.Value}, err
}

func (a *Auth) VerifyToken(ctx context.Context, token Token) (VerifyTokenResponse, error) {
	const op = "gRPC Auth VerifyToken"

	verifyToken, err := a.client.VerifyToken(ctx, &pb.Token{Value: token.Value})
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return VerifyTokenResponse{Ok: verifyToken.Ok}, err
}
