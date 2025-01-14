package exchange

import (
	"context"
	"fmt"
	pb "github.com/HennOgyrchik/proto-exchange/exchange"
)

func (e *Exchange) GetExchangeRates(ctx context.Context) (Rates, error) {
	const op = "gRPC Exchange GetExchangeRates"

	var result Rates

	rates, err := e.client.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		return result, fmt.Errorf("%s: %w", op, err)
	}

	return Rates{Rates: rates.Rates}, err
}

func (e *Exchange) GetExchangeRateForCurrency(ctx context.Context, in Currency) (Rate, error) {
	const op = "gRPC Exchange GetExchangeRateForCurrency"

	rate, err := e.client.GetExchangeRateForCurrency(ctx, &pb.CurrencyRequest{
		FromCurrency: in.FromCurrency,
		ToCurrency:   in.ToCurrency,
	})
	if err != nil {
		err = fmt.Errorf("%s: %w", op, err)
	}

	return Rate{
		FromCurrency: rate.FromCurrency,
		ToCurrency:   rate.ToCurrency,
		Rate:         rate.Rate,
	}, err
}
