package db

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func Test_TransferTx(t *testing.T) {
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)
	amount := int64(10)
	errChan := make(chan error)
	transferTxResultChan := make(chan *TransferTxResult)
	n := 5
	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			txResult, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAcountId:   fromAccount.AccountID,
				ToAcountId:     toAccount.AccountID,
				TransferAmount: amount,
			})

			errChan <- err
			transferTxResultChan <- txResult
		}()
	}

	fmt.Printf(
		">>Before Tx Balance account1=%d account2=%d \n",
		fromAccount.Balance,
		toAccount.Balance,
	)
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {

		if err := <-errChan; err != nil {
			t.Errorf("store.TransferTxFailed. got=%s", err)
			t.FailNow()
		}
		// check for transfer results
		txResult := <-transferTxResultChan

		if txResult.Transfer == nil {
			t.Errorf("txResult.Transfer is nil, transfer.id=%d", txResult.Transfer.ID)
		}

		if txResult.Transfer.FromAccountID != fromAccount.AccountID {
			t.Errorf(
				"txResult.Transfer.FromAccountID is not equal to fromAccount. expected=%d got=%d",
				fromAccount.AccountID,
				txResult.Transfer.FromAccountID,
			)
		}

		if txResult.Transfer.ToAccountID != toAccount.AccountID {
			t.Errorf("txResult.Transfer.ToAccountID is not equal to toAccount. expected=%d got=%d",
				fromAccount.AccountID,
				txResult.Transfer.ToAccountID)
		}

		if txResult.Transfer.Amount != amount {
			t.Errorf("txResult.Transfer.Amount is not equal to i, expected=%d got=%d",
				amount,
				txResult.Transfer.Amount)
		}

		if txResult.Transfer.ID == int64(0) {
			t.Errorf("txResult.Transfer.ID is 0")
		}

		if txResult.Transfer.CreatedAt.Time.IsZero() {
			t.Errorf("txResult.CreatedAt is zero")
		}

		_, err := testStore.GetTransfer(context.Background(), txResult.Transfer.ID)
		if err != nil {
			t.Errorf("error querying transaction with id=%d. got err=%v", txResult.Transfer.ID,
				err)
			t.FailNow()
		}

		// check entries
		// from entry
		fromEntry := txResult.FromEntry
		if fromEntry == nil {
			t.Error("txResult.FromEntry is nil")
			t.FailNow()
		}

		if fromEntry.AccountID != fromAccount.AccountID {
			t.Errorf("txResult.FromEntry.AccountId expected=%d, got=%d", fromAccount.AccountID,
				fromEntry.AccountID)
		}

		if fromEntry.Amount != -amount {
			t.Errorf("txResult.FromEntry.Amount expected=%d, got=%d", -amount,
				fromEntry.Amount)
		}

		if fromEntry.ID == int64(0) {
			t.Errorf("txResult.FromEntry.ID is 0")
		}
		if fromEntry.CreatedAt.Time.IsZero() {
			t.Errorf("txResult.CreatedAt is zero")
		}

		// to entry

		toEntry := txResult.ToEntry
		if toEntry == nil {
			t.Error("txResult.ToEntry is nil")
			t.FailNow()
		}

		if toEntry.AccountID != toAccount.AccountID {
			t.Errorf("txResult.ToEntry.AccountId expected=%d, got=%d", fromAccount.AccountID,
				toEntry.AccountID)
		}

		if toEntry.Amount != amount {
			t.Errorf("txResult.ToEntry.Amount expected=%d, got=%d", amount,
				toEntry.Amount)
		}

		if toEntry.ID == int64(0) {
			t.Errorf("txResult.ToEntry.ID is 0")
		}
		if toEntry.CreatedAt.Time.IsZero() {
			t.Errorf("txResult.CreatedAt is zero")
		}

		// check accounts

		fromAccountId := txResult.FromAccountId

		if reflect.ValueOf(fromAccountId).IsZero() {
			t.Error("txResult.FromAccountId is zero")
		}

		if fromAccountId.AccountID != fromAccount.AccountID {
			t.Errorf("txResult.FromAccountID.AccountId expected=%d. got=%d", fromAccount.AccountID,
				fromAccountId.AccountID)
		}

		toAccountId := txResult.ToAccountId

		if reflect.ValueOf(toAccountId).IsZero() {
			t.Error("txResult.ToAccountId is zero")
		}

		if toAccountId.AccountID != toAccount.AccountID {
			t.Errorf("txResult.ToAccountId.AccountId expected=%d. got=%d", toAccount.AccountID,
				toAccountId.AccountID)
		}

		// check for balance

		fmt.Printf("Tx account1=%d account2=%d \n", fromAccountId.Balance, toAccountId.Balance)
		diff1 := fromAccount.Balance - fromAccountId.Balance
		diff2 := toAccountId.Balance - toAccount.Balance

		if diff1 != diff2 {
			t.Errorf("diff1 != diff2. got diff1=%d diff2=%d", diff1, diff2)
		}
		if diff1 < 0 {
			t.Errorf("diff1 is less than zero. got=%d", diff1)
		}
		if (diff1 % amount) != 0 {
			t.Errorf("diff mod amount is different from zero. got=%d", diff1%amount)
		}

		k := int(diff1 / amount)
		if !(k >= 1 && k <= n) {
			t.Errorf("k is out of boundaries. got=%d", k)
		}

		ok, _ := existed[k]
		if ok {
			t.Errorf("k=%d in existed \n", k)
		}

		existed[k] = true

	}

	// check for final balance in accounts

	account1, _ := testStore.GetAccount(context.Background(), fromAccount.AccountID)
	account2, _ := testStore.GetAccount(context.Background(), toAccount.AccountID)

	fmt.Printf(">> Final account1=%d account2=%d", account1.Balance, account2.Balance)
	if int(account1.Balance) != (int(fromAccount.Balance) - n*int(amount)) {
		t.Errorf(
			"account1 balance is not correct expected=%d. got=%d \n",
			int(fromAccount.Balance)-n*int(amount),
			account1.Balance,
		)
	}
	if int(account2.Balance) != (int(toAccount.Balance) + n*int(amount)) {
		t.Errorf(
			"account2 balance is not correct expected=%d. got=%d \n",
			int(toAccount.Balance)+n*int(amount),
			account2.Balance,
		)
	}
}

func Test_TransferTxDeadlock(t *testing.T) {
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)
	amount := int64(10)
	errChan := make(chan error)
	n := 10
	for i := 0; i < n; i++ {

		outboundAccount := fromAccount.AccountID
		inboundAccount := toAccount.AccountID

		if i%2 == 0 {
			outboundAccount, inboundAccount = inboundAccount, outboundAccount
		}

		go func() {
			ctx := context.Background()
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAcountId:   outboundAccount,
				ToAcountId:     inboundAccount,
				TransferAmount: amount,
			})

			errChan <- err
		}()
	}

	fmt.Printf(
		">>Before Tx Balance account1=%d account2=%d \n",
		fromAccount.Balance,
		toAccount.Balance,
	)
	for i := 0; i < n; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("store.TransferTxFailed. got=%s", err)
			t.FailNow()
		}
		// check for transfer results
	}

	// check for final balance in accounts

	account1, _ := testStore.GetAccount(context.Background(), fromAccount.AccountID)
	account2, _ := testStore.GetAccount(context.Background(), toAccount.AccountID)

	fmt.Printf(">> Final account1=%d account2=%d", account1.Balance, account2.Balance)
	if int(account1.Balance) != int(fromAccount.Balance) {
		t.Errorf(
			"account1 balance is not correct expected=%d. got=%d \n",
			int(fromAccount.Balance)-n*int(amount),
			account1.Balance,
		)
	}
	if int(account2.Balance) != int(toAccount.Balance) {
		t.Errorf(
			"account2 balance is not correct expected=%d. got=%d \n",
			int(toAccount.Balance)+n*int(amount),
			account2.Balance,
		)
	}
}
