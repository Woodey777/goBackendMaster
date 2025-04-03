package db_sqlc

import (
	"context"
	"database/sql"
	"myBank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	trans1 := CreateRandomTransfer(t)

	trans2, err := testQueries.GetTransfer(context.Background(), trans1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, trans1)
	require.Equal(t, trans1.ID, trans2.ID)
	require.Equal(t, trans1.FromAccountID, trans2.FromAccountID)
	require.Equal(t, trans1.ToAccountID, trans2.ToAccountID)
	require.Equal(t, trans1.Amount, trans2.Amount)
	require.Equal(t, trans1.CreatedAt, trans2.CreatedAt)
}

func TestDeleteTransfer(t *testing.T) {
	trans1 := CreateRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), trans1.ID)
	require.NoError(t, err)
	
	trans2, err := testQueries.GetTransfer(context.Background(), trans1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, trans2)
}

func TestListTransfers(t *testing.T) {
	for range 10 { 
		CreateRandomTransfer(t)
	}

	args := ListTransfersParams{
		Limit: 5,
		Offset: 5,
	}

	entries, err := testQueries.ListTransfers(context.Background(), args)
	
	require.NoError(t, err)
	require.Len(t, entries, int(args.Limit))

	for _, acc := range entries {
		require.NotEmpty(t, acc)
	}
}

func CreateRandomTransfer(t *testing.T) Transfer {
	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)

	args := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID: acc2.ID,
		Amount:  util.RandomMoney(),
	}

	trans, err := testQueries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, trans)
	require.NotZero(t, trans.ID)
	require.NotZero(t, trans.CreatedAt)
	require.Equal(t, args.FromAccountID, trans.FromAccountID)
	require.Equal(t, args.ToAccountID, trans.ToAccountID)
	require.Equal(t, args.Amount, trans.Amount)

	return trans
}