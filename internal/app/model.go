package app

import (
	"context"
	"gw-currency-wallet/internal/grpcClient/exchange"
	"gw-currency-wallet/internal/storages"
)

type App struct {
	ctx       context.Context
	storage   storages.Storage
	exchanger exchange.Exchanger
}
