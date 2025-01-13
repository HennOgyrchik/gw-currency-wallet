package in_mem

import (
	"context"
	"sync"
)

type InMem struct {
	ctx  context.Context
	data sync.Map
}
