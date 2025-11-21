package auth

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type AdminStore interface {
	IsAdmin(email string) (bool, error)
	AddAdmin(email string) error
	RemoveAdmin(email string) error
}

type adminDB struct {
	db *sql.DB
}

func NewAdminDB(dbURL string) (AdminStore, error) {
	// Database handle is created for user store
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &adminDB{db: db}, nil
}

// IsAdmin checks if email corresponds to any admin record in the database
func (a *adminDB) IsAdmin(email string) (bool, error) {
	var count int
	err := a.db.QueryRow("SELECT COUNT(*) FROM admins WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	// Return true if count > 0
	return count > 0, nil
}

// AddAdmin adds the email address to the admins database table.
// This will allow users with this email to access the website
func (a *adminDB) AddAdmin(email string) error {
	exists, err := a.IsAdmin(email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user is already an admin: %s\n", email)
	}
	_, err = a.db.Exec("INSERT INTO admins(email) VALUES (?)", email)
	return err
}

// RemoveAdmin will remove this email address from the admins database table.
// If email doesn't exist in the database, function returns nil.
func (a *adminDB) RemoveAdmin(email string) error {
	_, err := a.db.Exec("DELETE FROM admins WHERE email = ?", email)
	if err != nil {
		return err
	}
	return nil
}
