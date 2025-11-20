package transaction

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type mariaDBStore struct {
	db *sql.DB
}

func NewMariaDBStore(dbURL string) (TransactionStore, error) {
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &mariaDBStore{db: db}, nil
}

func (s *mariaDBStore) Close() error {
	return s.db.Close()
}

func (s *mariaDBStore) execute(userID, amount int) (int, error) {
	result, err := s.db.Exec("INSERT INTO transactions (userId, amount) VALUES (?, ?)", userID, amount)
	if err != nil {
		return 0, err
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(transactionID), nil
}
