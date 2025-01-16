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
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Cash struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

type ExchangeRequest struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float32 `json:"amount"`
}

type ErrResponseJSON struct {
	Error string `json:"error"`
}

type MessageResponseJSON struct {
	Message string `json:"message"`
}

type TokenResponseJSON struct {
	Token string `json:"token"`
}

type NewBalanceResponseJSON struct {
	Message    string           `json:"message"`
	NewBalance storages.Balance `json:"new_balance"`
}

type ExchangeResponseJSON struct {
	Message        string           `json:"message"`
	ExchangeAmount float32          `json:"exchange_amount"`
	NewBalance     storages.Balance `json:"new_balance"`
}
