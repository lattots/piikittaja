package user

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/lattots/piikittaja/pkg/transaction"
)

type User struct {
	ID       int
	Username string
	Balance  int
	db       *sql.DB
}

// NewUser returns pointer to user instance. Function will fetch all
// user information from database or create a new user entry to database.
func NewUser(id int, username string) (*User, error) {
	// Database URL is read from environment variables.
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_APP environment variable is not set")
	}

	// Database handle is created for user.
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// User struct is created.
	user := &User{
		ID:       id,
		Username: username,
		db:       db,
	}

	// Database is checked for existing user.
	exists, err := user.Exists()
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if !exists {
		// If user doesn't already exist, new database entry is created.
		if err := user.initUser(); err != nil {
			return nil, fmt.Errorf("failed to initialize user: %w", err)
		}
	}

	// Method returns a pointer to user and a nil error.
	return user, nil
}

// initUser creates a new user entry to the database.
func (u *User) initUser() error {
	// User's username is inserted to database.
	// Balance gets a default value of 0 for all new user's.
	_, err := u.db.Exec("INSERT INTO users (id, username) VALUES (?, ?)", u.ID, u.Username)
	return err
}

// Exists checks if user already exists in the user database.
func (u *User) Exists() (bool, error) {
	var count int
	err := u.db.QueryRow("SELECT COUNT(*) FROM users WHERE id=?", u.ID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return count > 0, nil
}

// GetUsername returns the username of the user.
func (u *User) GetUsername() (string, error) {
	if u.Username != "" {
		return u.Username, nil
	}
	row := u.db.QueryRow("SELECT username FROM users WHERE id=?", u.ID)

	var username string
	err := row.Scan(&username)

	// If user is not found, function errors.
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("no user found with ID %d", u.ID)
	}

	// If error is other than not found, function returns the error as is.
	if err != nil {
		return "", fmt.Errorf("error fetching username from db: %w", err)
	}

	// User structs username is set.
	u.Username = username
	return username, nil
}

// GetBalance returns the current tab of the user.
func (u *User) GetBalance() (int, error) {
	row := u.db.QueryRow("SELECT balance FROM users WHERE id=?", u.ID)

	var balance int
	// Balance is scanned
	err := row.Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("no user found with id %d", u.ID)
	} else if err != nil {
		return 0, fmt.Errorf("failed to fetch user balance: %w", err)
	}

	return balance, nil
}

// AddToTab adds the amount to the user's tab. Returns the transaction ID and an error.
func (u *User) AddToTab(amount int) (int, error) {
	// Tabs are represented as negative balances in the database.
	// When adding to the tab, the amount is subtracted from current balance.
	return transaction.New(u.db, u.ID, -amount)
}

// PayBackTab subtracts the amount from user's tab. Returns the transaction ID and an error.
func (u *User) PayBackTab(amount int) (int, error) {
	// Tabs are represented as negative balances in the database.
	// When paying back the tab, the amount is added to current balance.
	return transaction.New(u.db, u.ID, amount)
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	row := db.QueryRow("SELECT id FROM users WHERE username=?", username)

	var user User
	user.Username = username
	user.db = db
	if err := row.Scan(&user.ID); err != nil {
		return nil, err
	}

	return &user, nil
}
