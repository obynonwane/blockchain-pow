// Package public maintains the group of handlers for public access.
package public

import (
	"context"
	"net/http"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
	"github.com/ardanlabs/blockchain/foundation/blockchain/state"
	"github.com/ardanlabs/blockchain/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of bar ledger endpoints.
type Handlers struct {
	Log   *zap.SugaredLogger
	State *state.State
}

// Genesis returns the genesis information
func (h Handlers) Genesis(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	gen := h.State.Genesis()
	return web.Respond(ctx, w, gen, http.StatusOK)
}

func (h *Handlers) Accounts(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// extract param from url
	accountStr := web.Param(r, "account")

	// declare a map
	var accounts map[database.AccountID]database.Account
	// switch statement for condition
	switch accountStr {
	case "":
		accounts = h.State.Accounts()
	default:
		accountID, err := database.ToAccountID(accountStr)
		if err != nil {
			return err
		}

		account, err := h.State.QueryAccount(accountID)
		if err != nil {
			return err
		}

		// construct a map
		accounts = map[database.AccountID]database.Account{accountID: account}
	}

	return web.Respond(ctx, w, accounts, http.StatusOK)
}
