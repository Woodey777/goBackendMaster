package db_sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
	*Queries
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Trasfer     Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *Store) TransferMoneyTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// create transfer
		result.Trasfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		// create entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		
		// update balances
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddAmountToAccount(ctx, AddAmountToAccountParams{
				Amount: -arg.Amount,
				ID: 		arg.FromAccountID,
			})
			if err != nil {
				return err
			}
	
			result.ToAccount, err = q.AddAmountToAccount(ctx, AddAmountToAccountParams{
				Amount: arg.Amount,
				ID:     arg.ToAccountID,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.AddAmountToAccount(ctx, AddAmountToAccountParams{
				Amount: arg.Amount,
				ID:     arg.ToAccountID,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddAmountToAccount(ctx, AddAmountToAccountParams{
				Amount: -arg.Amount,
				ID: 		arg.FromAccountID,
			})
			if err != nil {
				return err
			}
		}
		

		return nil
	})
	
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *Store) execTx(ctx context.Context, f func(q *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = f(q) 
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}
