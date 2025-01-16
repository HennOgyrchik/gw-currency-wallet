package in_mem

import (
	"sync"
	"time"
)

type InMem struct {
	data     sync.Map
	lifetime time.Duration
}
