package in_mem

import (
	"context"
	"gw-currency-wallet/internal/grpcClient/exchange"
	"sync"
	"time"
)

const recordKey = 1

func New(ctx context.Context, lifetime time.Duration) *InMem {

	ctxInMem, _ := context.WithCancel(ctx)

	inMem := InMem{
		data: sync.Map{},
		ctx:  ctxInMem,
	}

	go inMem.runCleaner(lifetime)

	return &inMem
}

func (i *InMem) Close() {
	i.ctx.Done()
}

func (i *InMem) GetRates() (exchange.Rates, bool) {
	var data exchange.Rates

	v, ok := i.data.Load(recordKey)
	if ok {
		data = v.(exchange.Rates)
	}

	return data, ok
}

func (i *InMem) RefreshRates(new exchange.Rates) {

	i.data.Store(recordKey, new)

}

func (i *InMem) delete(key any) {
	i.data.Delete(key)
}

func (i *InMem) runCleaner(timeout time.Duration) {
	for {
		select {
		case <-i.ctx.Done():
			return
		case <-time.After(timeout):
			i.delete(recordKey)
		}
	}
}
