package transaction

import "database/sql"

func New(db *sql.DB, userId int, amount int) error {
	// User's username is inserted to database.
	// Tab gets a default value of 0 for all new user's.
	_, err := db.Exec("INSERT INTO transactions (userId, amount) VALUES (?, ?)", userId, amount)
	if err != nil {
		return err
	}

	return nil
}
