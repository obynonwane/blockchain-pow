// Package mempool maintains the mempool for the blockchain
package mempool

import (
	"sync"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
)

// Mempool represents a cache of transactions organised by account:nonce
type Mempool struct {
	mu       sync.RWMutex
	pool     map[string]database.BlockTx
	selectFn selector.Func
}

// New constructs a new mempool using the default sort strategy
func New() (*Mempool, error) {
	return NewWithStrategy(selector.StrategyTip)
}
