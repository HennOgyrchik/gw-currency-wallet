package app

import (
	"context"
	"gw-currency-wallet/internal/cache"
	"gw-currency-wallet/internal/grpcClient/auth"
	"gw-currency-wallet/internal/grpcClient/exchange"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/logs"
)

type App struct {
	ctx        context.Context
	storage    storages.Storage
	cache      cache.Cache
	exchanger  exchange.Exchanger
	authorizer auth.Authorizer
	logger     *logs.Log
}

type User struct {
	Username string
	Password string
	Email    string
}

type Credentials struct {
	Username string
	Password string
}

type Cash struct {
	Amount   float32
	Currency string
}

type ExchangeRequest struct {
	FromCurrency string
	ToCurrency   string
	Amount       float32
}
