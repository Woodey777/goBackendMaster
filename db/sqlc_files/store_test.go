package db_sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferMoneyTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)

	n := 5
	resCh := make(chan TransferTxResult, n)
	errCh := make(chan error, n)
	amount := int64(100)
	// execute tx transfers
	for range n {
		go func () {
			transferResult, err := store.TransferMoneyTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount: 			 amount,
			})

			resCh <- transferResult
			errCh <- err
		}()
	}
	
	existed := make(map[int]bool, n)
	for range n {
		err := <-errCh
		require.NoError(t, err)
		
		transRes := <-resCh
		require.NotEmpty(t, transRes)
		
		// check transfer
		transfer := transRes.Trasfer
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entry
		fromEntry := transRes.FromEntry
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		toEntry := transRes.ToEntry
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		// check balance
		fromAccount := transRes.FromAccount
		require.NotZero(t, fromAccount.ID)
		require.NotZero(t, fromAccount.CreatedAt)
		require.Equal(t, acc1.Owner,fromAccount.Owner)
		require.Equal(t, acc1.Currency,fromAccount.Currency)

		toAccount := transRes.ToAccount
		require.NotZero(t, toAccount.ID)
		require.NotZero(t, toAccount.CreatedAt)
		require.Equal(t, acc2.Owner,toAccount.Owner)
		require.Equal(t, acc2.Currency,toAccount.Currency)


		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 >= 0)
		require.True(t, diff1 % amount == 0)

		k := int(diff1/amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	getAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	
	getAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	balanceDiff := amount * int64(n)
	require.Equal(t, acc1.Balance - balanceDiff, getAcc1.Balance)
	require.Equal(t, acc2.Balance + balanceDiff, getAcc2.Balance)
}

func TestTransferMoneyTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)

	n := 10
	errCh := make(chan error, n)
	amount := int64(100)
	// execute tx transfers
	for i := 0; i < n; i++ {
		fromAccountID := acc1.ID
		toAccountID := acc2.ID

		if i % 2 == 0 {
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}

		go func () {

			_, err := store.TransferMoneyTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount: 			 amount,
			})

			errCh <- err
		}()
	}
	
	for range n {
		err := <-errCh
		require.NoError(t, err)
	}

	getAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	
	getAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, getAcc1.Balance)
	require.Equal(t, acc2.Balance, getAcc2.Balance)
}
