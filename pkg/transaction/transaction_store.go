package transaction

import (
	"time"

	"github.com/lattots/piikittaja/pkg/models"
)

type TransactionStore interface {
	execute(userID, amount int) (transactionID int, err error)
	getTransactions(endDate time.Time, window time.Duration, traType string) ([]*models.Transaction, error)
	getUserTransactions(userId, quantity int) ([]*models.Transaction, error)
	Close() error
}
