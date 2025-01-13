package storages

import (
	"context"
)

type Storage interface {
	GetBalance(context.Context, string) (Balance, error)
	NewWallet(context.Context, string) error
	UpdateWallet(context.Context, string, Balance) error
}

type Balance struct {
	USD,
	RUB,
	EUR float32
}
