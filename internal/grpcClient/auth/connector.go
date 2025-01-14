package auth

import (
	"fmt"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(grpcServerURL string) *Auth {
	return &Auth{
		url:    grpcServerURL,
		conn:   nil,
		client: nil,
	}
}

func (a *Auth) Run() error {
	const op = "gRPC Auth New"

	conn, err := grpc.NewClient(a.url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.conn = conn
	a.client = pb.NewAuthorizationClient(conn)

	return nil
}

func (a *Auth) Stop() error {
	const op = "gRPC Auth Stop"

	if err := a.conn.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
