package cache

import "gw-currency-wallet/internal/grpcClient/exchange"

type Cache interface {
	Close()
	GetRates() (exchange.Rates, bool)
	RefreshRates(rates exchange.Rates)
}
