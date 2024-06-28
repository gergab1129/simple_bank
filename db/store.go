package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
			return fmt.Errorf("tx err: %s, rb error: %s", err, rollBackErr)
		}

		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAcountId   int64 `json:"fromAccountId"`
	ToAcountId     int64 `json:"toAccountId"`
	TransferAmount int64 `json:"transferAmount"`
}

type TransferTxResult struct {
	Transfer      *Transfers `json:"transferId"`
	FromAccountId Accounts   `json:"fromAccountId"`
	ToAccountId   Accounts   `json:"ToAccountId"`
	FromEntry     *Entries   `json:"fromEntry"`
	ToEntry       *Entries   `json:"toEntry"`
}

var TxKey = struct{}{}

// TrasferTx creates a money trasnfer between accounts.
// perfoms the following operations create a transfer record, create accounts entries and update account balances
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error) {
	txResults := &TransferTxResult{}

	err := s.execTx(ctx, func(q *Queries) error {
		TransferResult, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAcountId,
			ToAccountID:   arg.ToAcountId,
			Amount:        arg.TransferAmount,
		})
		if err != nil {
			return err
		}

		txResults.Transfer = &TransferResult

		FromEntryResults, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAcountId,
			Amount:    -arg.TransferAmount,
		})
		if err != nil {
			return err
		}
		txResults.FromEntry = &FromEntryResults

		ToEntryResult, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAcountId,
			Amount:    arg.TransferAmount,
		})
		if err != nil {
			return err
		}
		txResults.ToEntry = &ToEntryResult

		if arg.FromAcountId < arg.ToAcountId {
			txResults.FromAccountId, txResults.ToAccountId, err = AddBalance(ctx, q,
				AddBalanceParams{
					AccountID: arg.FromAcountId,
					Amount:    -arg.TransferAmount,
				},
				AddBalanceParams{
					AccountID: arg.ToAcountId,
					Amount:    arg.TransferAmount,
				},
			)
			if err != nil {
				return err
			}
		} else {
			txResults.FromAccountId, txResults.ToAccountId, err = AddBalance(ctx, q,
				AddBalanceParams{
					AccountID: arg.ToAcountId,
					Amount:    arg.TransferAmount,
				},
				AddBalanceParams{
					AccountID: arg.FromAcountId,
					Amount:    -arg.TransferAmount,
				},
			)
			if err != nil {
				return err
			}

		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return txResults, nil
}

func AddBalance(ctx context.Context,
	q *Queries,
	args ...AddBalanceParams,
) (Accounts, Accounts, error) {
	account1, err := q.AddBalance(ctx, args[0])
	if err != nil {
		return Accounts{}, Accounts{}, err
	}

	account2, err := q.AddBalance(ctx, args[1])
	if err != nil {
		return Accounts{}, Accounts{}, err
	}

	return account1, account2, nil
}
