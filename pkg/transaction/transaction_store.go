package transaction

type TransactionStore interface {
	execute(userID, amount int) (transactionID int, err error)
	Close() error
}
