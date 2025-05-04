package state

import "github.com/ardanlabs/blockchain/foundation/blockchain/database"

func (s *State) QueryAccount(account database.AccountID) (database.Account, error) {
	return s.db.Query(account)
}
