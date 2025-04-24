// Package mempool maintains the mempool for the blockchain
package mempool

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
)

// Mempool represents a cache of transactions organised by account:nonce
type Mempool struct {
	mu   sync.RWMutex
	pool map[string]database.BlockTx // string here which is the key is the concate of from addresss and nonce
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

	// Ethereum requires a 10% bump in the tip to replace an existing
	// transaction in the mempool and so do we. We want to limit users
	// from this sort of behavior.
	if etx, exists := mp.pool[key]; exists {
		if tx.Tip < uint64(math.Round(float64(etx.Tip)*1.10)) {
			return errors.New("replacing a transaction requires a 10% bump in the tip")
		}
	}

	// update the tx in the mempool
	// the key is the concate with (:) of from address and nonce
	mp.pool[key] = tx
	return nil
}

// Delete removed a transaction from the mempool
func (mp *Mempool) Delete(tx database.BlockTx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	key, err := mapKey(tx)
	if err != nil {
		return err
	}

	// delete item from the pool
	delete(mp.pool, key)

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
