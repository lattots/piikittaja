package transaction

import "github.com/lattots/piikittaja/pkg/models"

type TransactionStore interface {
	execute(userID, amount int) (transactionID int, err error)
	getTransactions(userId, quantity int) ([]*models.Transaction, error)
	Close() error
}
