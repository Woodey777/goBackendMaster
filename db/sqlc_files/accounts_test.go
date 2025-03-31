package db_sqlc

import (
	"testing"
	"util"
)

func TestCreateAccount(t *testing.T) {
	params := CreateAccountParams{
		Owner:		"John Doe",
		Balance:	1000,
		Currency: "USD",
	}
}
