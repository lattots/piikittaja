package transaction

import "github.com/lattots/piikittaja/pkg/models"

type mockStore struct {
	transactions map[int]mockEntry
	nextID       int
}

type mockEntry struct {
	userID int
	amount int
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
	s.transactions[transactionID] = mockEntry{userID: userID, amount: amount}

	return transactionID, nil
}

func (s *mockStore) getTransactions(userID, quantity int) ([]*models.Transaction, error) {
	var res []*models.Transaction

	for _, t := range s.transactions {
		if t.userID == userID {
			resTransaction := &models.Transaction{
				Amount: t.amount,
			}
			if resTransaction.Amount < 0 {
				resTransaction.Type = "withdraw"
				resTransaction.Amount = -resTransaction.Amount
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
