package db

import (
	"context"
	"testing"
)

func Test_TransferTx(t *testing.T) {
	fromAccount := createTestAccount(t).AccountID
	toAccount := createTestAccount(t).AccountID
	amount := int64(10)
	errChan := make(chan error)
	transferTxResultChan := make(chan *TransferTxResult)
	for i := 0; i < 10; i++ {
		go func() {
			txResult, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAcountId:   fromAccount,
				ToAcountId:     toAccount,
				TransferAmount: amount,
			})

			errChan <- err
			transferTxResultChan <- txResult
		}()
	}

	for i := 0; i < 10; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("store.TransferTxFailed. got=%s", <-errChan)
			t.FailNow()
		}
		// check for transfer results
		txResult := <-transferTxResultChan

		if txResult.Transfer == nil {
			t.Errorf("txResult.Transfer is nil, transfer.id=%d", txResult.Transfer.ID)
		}

		if txResult.Transfer.FromAccountID != fromAccount {
			t.Errorf(
				"txResult.Transfer.FromAccountID is not equal to fromAccount. expected=%d got=%d",
				fromAccount,
				txResult.Transfer.FromAccountID,
			)
		}

		if txResult.Transfer.ToAccountID != toAccount {
			t.Errorf("txResult.Transfer.ToAccountID is not equal to toAccount. expected=%d got=%d",
				fromAccount,
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

		if fromEntry.AccountID != fromAccount {
			t.Errorf("txResult.FromEntry.AccountId expected=%d, got=%d", fromAccount,
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

		if toEntry.AccountID != toAccount {
			t.Errorf("txResult.ToEntry.AccountId expected=%d, got=%d", fromAccount,
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

	}
}
