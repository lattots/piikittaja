package transaction

import "database/sql"

// New executes a new transaction. Returns the transaction ID and an error.
func New(db *sql.DB, userId int, amount int) (int, error) {
	// User's username is inserted to database.
	// Tab gets a default value of 0 for all new user's.
	result, err := db.Exec("INSERT INTO transactions (userId, amount) VALUES (?, ?)", userId, amount)
	if err != nil {
		return 0, err
	}

	transactionId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(transactionId), nil
}
