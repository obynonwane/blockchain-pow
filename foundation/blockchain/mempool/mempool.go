// Package mempool maintains the mempool for the blockchain
package mempool

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
)

// Mempool represents a cache of transactions organised by account:nonce
type Mempool struct {
	mu   sync.RWMutex
	pool map[string]database.BlockTx
}

// New constructs a new mempool using the default sort strategy
func New() (*Mempool, error) {
	return NewWithStrategy()
}

// NewWithStrategy constructs a new mempool with specified sort strategy
func NewWithStrategy() (*Mempool, error) {
	mp := Mempool{
		pool: make(map[string]database.BlockTx),
	}
	return &mp, nil
}

// Count returns the current number of transactions in the pool.
func (mp *Mempool) Count() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return len(mp.pool)
}

// Upsert adds or replace a transaction from th mempool
func (mp *Mempool) Upsert(tx database.BlockTx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// CORE NOTE: Different blockchains have different algorithms to limit the
	// size of the mempool. Some limit based on the amount of memory being
	// consumed and some may limit based on the number of transaction. If a limit
	// is met, then either the transaction that has the least return on investment
	// or the oldest will be dropped from the pool to make room for new the transaction.

	// For now my blockchain is not imposing any limits
	key, err := mapKey(tx)
	if err != nil {
		return err
	}

	return nil
}

// ========================================================================================

// mapKey is used to generate the map key.
func mapKey(tx database.BlockTx) (string, error) {
	return fmt.Sprintf("%s:%d", tx.FromID, tx.Nonce), nil
}

// accountFromMapKey extracts the account information from the mapKey
func accountFromMapKey(key string) database.AccountID {
	// splits key string and convert to account type returning th e0 indexed
	return database.AccountID(strings.Split(key, ":")[0])
}
