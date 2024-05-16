package user

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

type User struct {
	id       int
	username string
	db       *sql.DB
}

// NewUser returns pointer to user instance. Function will fetch all
// user information from database or create a new user entry to database.
func NewUser(username string) (*User, error) {
	// Database handle is created for user.
	db, err := sql.Open("mysql", os.Getenv("DATABASE_APP"))
	if err != nil {
		return nil, err
	}

	// User struct is created.
	user := User{
		username: username,
		db:       db,
	}

	// Database is checked for existing user.
	exists, err := user.Exists()
	if err != nil {
		return nil, err
	}

	if !exists {
		// If user doesn't already exist, new database entry is created.
		err = user.initUser()
		if err != nil {
			return nil, err
		}
	} else {
		// If user exists, user data is fetched from database.
		err := user.getIdByUsername()
		if err != nil {
			return nil, err
		}
	}

	// Method returns a pointer to user and a nil error.
	return &user, nil
}

// initUser creates a new user entry to the database. Created user's ID is also added to the user's ID field.
func (u *User) initUser() error {
	// User's username is inserted to database.
	// Tab gets a default value of 0 for all new user's.
	result, err := u.db.Exec("INSERT INTO users (username, tab) VALUES (?, ?)", u.username, 0)
	if err != nil {
		return err
	}

	// Automatically created ID is read from the result.
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// User structs ID is set.
	u.id = int(id)
	return nil
}

// Exists checks if user already exists in the user database.
func (u *User) Exists() (bool, error) {
	var count int
	err := u.db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", u.username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetId returns the ID of the user.
func (u *User) GetId() (int, error) {
	if u.id != 0 {
		return u.id, nil
	}
	row := u.db.QueryRow("SELECT id FROM users WHERE username=?", u.username)

	var id int
	err := row.Scan(&id)

	// If user is not found, function errors.
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("no user found with username %s", u.username)
	}

	// If error is other than not found, function returns the error as is.
	if err != nil {
		return 0, err
	}

	// User structs id is set.
	u.id = id
	return id, nil
}

// GetUsername returns the username of the user.
func (u *User) GetUsername() (string, error) {
	if u.username != "" {
		return u.username, nil
	}
	row := u.db.QueryRow("SELECT username FROM users WHERE id=?", u.id)

	var username string
	err := row.Scan(&username)

	// If user is not found, function errors.
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("no user found with ID %d", u.id)
	}

	// If error is other than not found, function returns the error as is.
	if err != nil {
		return "", err
	}

	// User structs username is set.
	u.username = username
	return username, nil
}

// GetTab returns the current tab of the user.
func (u *User) GetTab() (int, error) {
	row := u.db.QueryRow("SELECT tab FROM users WHERE id=?", u.id)

	var tab int
	err := row.Scan(&tab)

	// If user is not found, function errors.
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("no user found with id %d", u.id)
	}

	// If error is other than not found, function returns the error as is.
	if err != nil {
		return 0, err
	}

	return tab, nil
}

// AddToTab adds the amount to the user's tab.
func (u *User) AddToTab(amount int) error {
	_, err := u.db.Exec("UPDATE users SET tab = tab + ? WHERE id = ?", amount, u.id)
	return err
}

// getIdByUsername gets user ID from the database and adds it to user structs id field.
func (u *User) getIdByUsername() error {
	row := u.db.QueryRow("SELECT id FROM users WHERE username=?", u.username)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}

	// User structs ID is set.
	u.id = id

	return nil
}
