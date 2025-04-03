package db_sqlc

import (
	"context"
	"database/sql"
	"myBank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := CreateRandomAccount(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.CreatedAt, acc2.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := CreateRandomAccount(t)

	args := UpdateAccountParams{
		ID: acc1.ID,
		Balance: util.RandomMoney(),
	}

	updatedAcc, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	require.Equal(t, acc1.ID, updatedAcc.ID)
	require.Equal(t, acc1.Owner, updatedAcc.Owner)
	require.Equal(t, args.Balance, updatedAcc.Balance)
	require.Equal(t, acc1.Currency, updatedAcc.Currency)
	require.Equal(t, acc1.CreatedAt, updatedAcc.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := CreateRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, acc2)
}

func TestListAccounts(t *testing.T) {
	for range 10 { 
		CreateRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	
	require.NoError(t, err)
	require.Len(t, accounts, int(args.Limit))

	for _, acc := range accounts {
		require.NotEmpty(t, acc)
	}
}

func CreateRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	require.Equal(t, args.Owner, acc.Owner)
	require.Equal(t, args.Balance, acc.Balance)
	require.Equal(t, args.Currency, acc.Currency)

	return acc
}