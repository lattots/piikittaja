package transaction

import (
	"errors"
	"fmt"

	"github.com/lattots/piikittaja/pkg/models"
)

type TransactionHandler interface {
	Withdraw(user *models.User, amount int) (int, error)
	Deposit(user *models.User, amount int) (int, error)
	GetTransactions(user *models.User, quantity int) ([]*models.Transaction, error)
}

var ErrNotEnoughBalance = errors.New("user doesn't have enough balance to withdraw the amount")
var ErrInvalidAmount = errors.New("entered amount is invalid. amount must be greater than 0")

// canWithdraw is a helper function to check if user can withdraw specified amount
func canWithdraw(user *models.User, amount int) bool {
	const debtThreshold = 10_00 // = 10 €
	return user.Balance+debtThreshold >= amount || user.Username == "maanmittauskilta"
}

type transactionHandler struct {
	store TransactionStore
}

func NewTransactionHandler(store TransactionStore) TransactionHandler {
	return &transactionHandler{store: store}
}

func (h *transactionHandler) Withdraw(user *models.User, amount int) (int, error) {
	if amount <= 0 {
		return 0, ErrInvalidAmount
	}
	if !canWithdraw(user, amount) {
		return 0, ErrNotEnoughBalance
	}

	id, err := h.store.execute(user.ID, -amount)
	if err != nil {
		return 0, fmt.Errorf("error executing transaction: %w", err)
	}
	user.Balance -= amount

	return id, nil
}

func (h *transactionHandler) Deposit(user *models.User, amount int) (int, error) {
	if amount <= 0 {
		return 0, ErrInvalidAmount
	}

	id, err := h.store.execute(user.ID, amount)
	if err != nil {
		return 0, fmt.Errorf("error executing transaction: %w", err)
	}
	user.Balance += amount

	return id, nil
}

func (h *transactionHandler) GetTransactions(user *models.User, quantity int) ([]*models.Transaction, error) {
	return h.store.getTransactions(user.ID, quantity)
}
