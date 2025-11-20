package transaction_test

import (
	"errors"
	"os"
	"testing"

	"github.com/lattots/piikittaja/pkg/models"
	"github.com/lattots/piikittaja/pkg/transaction"
	userstore "github.com/lattots/piikittaja/pkg/user_store"
)

func TestMariaDBStore(t *testing.T) {
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		t.Fatal("DATABASE_APP environment variable is not set")
	}

	usrStore, err := userstore.NewMariaDBStore(dbURL)
	if err != nil {
		t.Fatalf("error creating user store: %s", err)
	}
	defer usrStore.Close()

	traStore, err := transaction.NewMariaDBStore(dbURL)
	if err != nil {
		t.Fatalf("error creating transaction store: %s", err)
	}
	defer traStore.Close()
	handler := transaction.NewTransactionHandler(traStore)

	const testUserID = 0
	user := models.NewUser(testUserID, "test user 0")

	err = usrStore.Insert(user)
	if err != nil {
		t.Fatalf("error inserting user to store: %s", err)
	}

	_, err = handler.Deposit(user, 100)
	if err != nil {
		t.Fatalf("error making deposit transaction: %s", err)
	}

	if user.Balance != 100 {
		t.Fatalf("bad user balance in memory: want 100 got %d", user.Balance)
	}

	_, err = handler.Withdraw(user, 50)
	if err != nil {
		t.Fatalf("error making withdraw transaction: %s", err)
	}

	if user.Balance != 50 {
		t.Fatalf("bad user balance in memory: want 50 got %d", user.Balance)
	}

	user, err = usrStore.GetByID(testUserID)
	if err != nil {
		t.Fatalf("error getting test user from store: %s", err)
	}

	if user.Balance != 50 {
		t.Fatalf("bad user balance in store: want 50 got %d", user.Balance)
	}

	_, err = handler.Withdraw(user, 61)
	if !errors.Is(err, transaction.ErrNotEnoughBalance) {
		t.Fatalf("wrong error for withdraw attempt: want ErrNotEnoughBalance got %s", err)
	}

	if user.Balance != 50 {
		t.Fatalf("bad user balance in store: want 50 got %d", user.Balance)
	}

	err = usrStore.Remove(testUserID)
	if err != nil {
		t.Fatalf("error removing inserted user: %s", err)
	}
}
