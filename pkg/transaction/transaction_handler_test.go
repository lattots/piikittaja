package transaction_test

import (
	"errors"
	"testing"

	"github.com/lattots/piikittaja/pkg/models"
	"github.com/lattots/piikittaja/pkg/transaction"
)

func TestValidTransactions(t *testing.T) {
	store := transaction.NewMockStore()
	handler := transaction.NewTransactionHandler(store)

	user := models.NewUser(0, "test user", "John", "Doe")

	_, err := handler.Deposit(user, 100_00)
	if err != nil {
		t.Fatalf("valid deposit failed: %s", err)
	}

	if user.Balance != 100_00 {
		t.Fatalf("user balance invalid: want 10000 got %d", user.Balance)
	}

	_, err = handler.Withdraw(user, 10_00)
	if err != nil {
		t.Fatalf("valid withdraw failed: %s", err)
	}

	if user.Balance != 90_00 {
		t.Fatalf("user balance invalid: want 9000 got %d", user.Balance)
	}

	_, err = handler.Withdraw(user, 100_00)
	if err != nil {
		t.Fatalf("valid withdraw failed: %s", err)
	}

	if user.Balance != -10_00 {
		t.Fatalf("user balance invalid: want -10_00 got %d", user.Balance)
	}
}

func TestInvalidTransactions(t *testing.T) {
	store := transaction.NewMockStore()
	handler := transaction.NewTransactionHandler(store)

	user := models.NewUser(0, "test user", "John", "Doe")

	_, err := handler.Deposit(user, -1)
	if !errors.Is(err, transaction.ErrInvalidAmount) {
		t.Fatalf("wrong error for invalid amount: want ErrInvalidAmount got %s", err)
	}

	if user.Balance != 0 {
		t.Fatalf("user balance invalid: want 0 got %d", user.Balance)
	}

	_, err = handler.Withdraw(user, -1)
	if !errors.Is(err, transaction.ErrInvalidAmount) {
		t.Fatalf("wrong error for invalid amount: want ErrInvalidAmount got %s", err)
	}

	if user.Balance != 0 {
		t.Fatalf("user balance invalid: want 0 got %d", user.Balance)
	}

	// Try to withdraw over the debt limit
	_, err = handler.Withdraw(user, 11_00) // 11 €
	if !errors.Is(err, transaction.ErrNotEnoughBalance) {
		t.Errorf("wrong error for insufficient funds: want ErrNotEnoughBalance got %s", err)
	}

	if user.Balance != 0 {
		t.Fatalf("user balance invalid: want 0 got %d", user.Balance)
	}
}
