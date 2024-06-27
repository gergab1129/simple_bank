package db

import (
	"context"
	"testing"

	"simple_bank/util"
)

func createTestAccount(t *testing.T) *Accounts {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	r, err := testQueries.CreateAccount(context.Background(), args)
	if err != nil {
		t.Fatal(err)
	}

	return &r
}

func Test_CreateAccount(t *testing.T) {
	account := CreateAccountParams{}

	account.Owner = "german"
	account.Balance = 100
	account.Currency = "USD"

	r, err := testQueries.CreateAccount(context.Background(), account)
	if err != nil {
		t.Errorf("error creating account, got=%s", err)
	}

	if r.Owner != account.Owner {
		t.Errorf(
			"r.Owner is not equal to account.Owner. Expected=%s got=%s",
			account.Owner,
			r.Owner,
		)
	}
	if r.Currency != account.Currency {
		t.Errorf(
			"r.Currency is not equal to account.Currency. Expected=%s got=%s",
			account.Currency,
			r.Currency,
		)
	}
	if r.Balance != account.Balance {
		t.Errorf(
			"r.Balance is not equal to account.Balance. Expected=%d got=%d",
			account.Balance,
			r.Balance,
		)
	}
}
