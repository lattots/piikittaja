package transaction

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lattots/piikittaja/pkg/models"
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

func (s *mariaDBStore) getTransactions(endDate time.Time, window time.Duration, traType string) ([]*models.Transaction, error) {
	startDate := endDate.Add(-window)

	var amountCondition string
	switch traType {
	case "deposit":
		amountCondition = "AND amount > 0"
	case "withdraw":
		amountCondition = "AND amount < 0"
	default:
		// If no type is specified, return all transactions
		amountCondition = ""
	}

	query := fmt.Sprintf("SELECT time, amount FROM transactions WHERE time >= ? AND time <= ? %s ORDER BY time DESC", amountCondition)

	rows, err := s.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		transaction := &models.Transaction{}
		var rawTime []uint8
		if err := rows.Scan(&rawTime, &transaction.Amount); err != nil {
			return nil, err
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(rawTime))
		if err != nil {
			return nil, err
		}
		transaction.IssuedAt = parsedTime

		// Since withdraw transactions are stored as negative amounts, transaction
		// type must be set accordingly and the amount must be inversed.
		if transaction.Amount < 0 {
			transaction.Type = "withdraw"
			transaction.Amount = -transaction.Amount
		} else {
			transaction.Type = "deposit"
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *mariaDBStore) getUserTransactions(userID, quantity int) ([]*models.Transaction, error) {
	rows, err := s.db.Query("SELECT time, amount FROM transactions WHERE userId=? ORDER BY time DESC LIMIT ?", userID, quantity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		transaction := new(models.Transaction)
		var rawTime []uint8
		err := rows.Scan(&rawTime, &transaction.Amount)
		if err != nil {
			return nil, err
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(rawTime))
		if err != nil {
			return nil, err
		}
		transaction.IssuedAt = parsedTime

		// Since withdraw transactions are stored as negative amounts, transaction
		// type must be set accordingly and the amount must be inversed.
		if transaction.Amount < 0 {
			transaction.Type = "withdraw"
			transaction.Amount = -transaction.Amount
		} else {
			transaction.Type = "deposit"
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
