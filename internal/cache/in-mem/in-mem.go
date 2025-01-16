package in_mem

import (
	"sync"
	"time"
)

func New(lifetime time.Duration) *InMem {

	inMem := InMem{
		data:     sync.Map{},
		lifetime: lifetime,
	}

	return &inMem
}

func (i *InMem) Get(key string) (any, bool) {
	return i.data.Load(key)
}

func (i *InMem) Set(key string, value any) {
	i.data.Store(key, value)

	time.AfterFunc(i.lifetime, func() {
		i.data.Delete(key)
	})

}
