package db_sqlc

import (
	"context"
	"database/sql"
	"myBank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	ent1 := CreateRandomEntry(t)

	ent2, err := testQueries.GetEntry(context.Background(), ent1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, ent1)
	require.Equal(t, ent1.ID, ent2.ID)
	require.Equal(t, ent1.AccountID, ent2.AccountID)
	require.Equal(t, ent1.Amount, ent2.Amount)
	require.Equal(t, ent1.CreatedAt, ent2.CreatedAt)
}

func TestDeleteEntry(t *testing.T) {
	ent1 := CreateRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), ent1.ID)
	require.NoError(t, err)
	
	ent2, err := testQueries.GetEntry(context.Background(), ent1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, ent2)
}

func TestListEntries(t *testing.T) {
	for range 10 { 
		CreateRandomEntry(t)
	}

	args := ListEntriesParams{
		Limit: 5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)
	
	require.NoError(t, err)
	require.Len(t, entries, int(args.Limit))

	for _, acc := range entries {
		require.NotEmpty(t, acc)
	}
}

func CreateRandomEntry(t *testing.T) Entry {
	acc := CreateRandomAccount(t)

	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:  util.RandomMoney(),
	}

	ent, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, ent)
	require.NotZero(t, ent.ID)
	require.NotZero(t, ent.CreatedAt)
	require.Equal(t, args.AccountID, ent.AccountID)
	require.Equal(t, args.Amount, ent.Amount)

	return ent
}