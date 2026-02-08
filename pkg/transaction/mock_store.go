package transaction

import (
	"time"

	"github.com/lattots/piikittaja/pkg/models"
)

type mockStore struct {
	transactions map[int]mockEntry
	nextID       int
}

type mockEntry struct {
	userID   int
	amount   int
	issuedAt time.Time
}

func NewMockStore() TransactionStore {
	transactions := make(map[int]mockEntry)
	return &mockStore{transactions: transactions, nextID: 0}
}

func (s *mockStore) Close() error {
	return nil
}

func (s *mockStore) execute(userID, amount int) (int, error) {
	transactionID := s.nextID
	s.nextID++
	newTra := mockEntry{
		userID:   userID,
		amount:   amount,
		issuedAt: time.Now(),
	}
	s.transactions[transactionID] = newTra

	return transactionID, nil
}

func (s *mockStore) getTransactions(endDate time.Time, window time.Duration, traType string) ([]*models.Transaction, error) {
	res := make([]*models.Transaction, 0)

	for _, t := range s.transactions {
		if t.issuedAt.Before(endDate) && t.issuedAt.After(endDate.Add(-window)) {
			resTransaction := &models.Transaction{
				Amount: t.amount,
			}
			if resTransaction.Amount < 0 {
				resTransaction.Type = "withdraw"
				resTransaction.Amount = -resTransaction.Amount
			} else {
				resTransaction.Type = "deposit"
			}
			res = append(res, resTransaction)
		}
	}

	return res, nil
}

func (s *mockStore) getUserTransactions(userID, quantity int) ([]*models.Transaction, error) {
	var res []*models.Transaction

	for _, t := range s.transactions {
		if t.userID == userID {
			resTransaction := &models.Transaction{
				Amount: t.amount,
			}
			if resTransaction.Amount < 0 {
				resTransaction.Type = "withdraw"
				resTransaction.Amount = -resTransaction.Amount
			} else {
				resTransaction.Type = "deposit"
			}
			res = append(res, resTransaction)

			// Only return up to quantity number of transactions
			if len(res) >= quantity {
				break
			}
		}
	}
	return res, nil
}
