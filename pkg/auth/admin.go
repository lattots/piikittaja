package auth

import (
	"database/sql"
)

type AdminStore interface {
	IsAdmin(email string) bool
	AddAdmin(email string) error
	RemoveAdmin(email string) error
}

type AdminDB struct {
	db *sql.DB
}

func NewAdminDB(db *sql.DB) *AdminDB {
	return &AdminDB{db: db}
}

// IsAdmin checks if email corresponds to any admin record in the database
func (a *AdminDB) IsAdmin(email string) bool {
	var exists bool
	err := a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM admins WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// AddAdmin adds the email address to the admins database table.
// This will allow users with this email to access the website
func (a *AdminDB) AddAdmin(email string) error {
	_, err := a.db.Exec("INSERT INTO admins(email) VALUES (?)", email)
	return err
}

// RemoveAdmin will remove this email address from the admins database table.
// If email doesn't exist in the database, function returns nil.
func (a *AdminDB) RemoveAdmin(email string) error {
	_, err := a.db.Exec("DELETE FROM admins WHERE email = ?", email)
	if err != nil {
		return err
	}
	return nil
}
