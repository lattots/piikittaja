package transaction

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

func (h *mockStore) execute(userID, amount int) (int, error) {
	transactionID := h.nextID
	h.nextID++
	h.transactions[transactionID] = mockEntry{userID: userID, amount: amount}

	return transactionID, nil
}
