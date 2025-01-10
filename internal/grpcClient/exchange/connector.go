package exchange

import (
	"fmt"
	pb "github.com/HennOgyrchik/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(grpcServerURL string) *Exchange {
	return &Exchange{
		url:    grpcServerURL,
		conn:   nil,
		client: nil,
	}
}

func (e *Exchange) Run() error {

	fmt.Println("123123123")

	const op = "gRPC Exchange New"

	conn, err := grpc.NewClient(e.url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	e.conn = conn
	e.client = pb.NewExchangeServiceClient(conn)

	fmt.Println("CLIENT = ", e.client, "FUNC = ", pb.NewExchangeServiceClient(conn))

	return nil
}

func (e *Exchange) Stop() error {
	const op = "gRPC Exchange Stop"

	if err := e.conn.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
