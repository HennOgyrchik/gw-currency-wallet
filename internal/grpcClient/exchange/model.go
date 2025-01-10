package exchange

import (
	"context"
	pb "github.com/HennOgyrchik/proto-exchange/exchange"
	"google.golang.org/grpc"
)

type Exchanger interface {
	GetExchangeRates(ctx context.Context) (Rates, error)
	GetExchangeRateForCurrency(ctx context.Context, in Currency) (Rate, error)
}

type Exchange struct {
	url    string
	conn   *grpc.ClientConn
	client pb.ExchangeServiceClient
}

type Rates struct {
	Rates map[string]float32
}

type Currency struct {
	FromCurrency string
	ToCurrency   string
}

type Rate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float32
}
