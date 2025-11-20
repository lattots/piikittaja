package userstore

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lattots/piikittaja/pkg/models"

	_ "github.com/go-sql-driver/mysql"
)

type mariaDBStore struct {
	db *sql.DB
}

func NewMariaDBStore(dbURL string) (UserStore, error) {
	// Database handle is created for user store
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

func (s *mariaDBStore) GetByID(id int) (*models.User, error) {
	row := s.db.QueryRow("SELECT username, balance FROM users WHERE id=?", id)

	user := &models.User{ID: id}
	if err := row.Scan(&user.Username, &user.Balance); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *mariaDBStore) GetByUsername(username string) (*models.User, error) {
	row := s.db.QueryRow("SELECT id, balance FROM users WHERE username=?", username)

	user := &models.User{Username: username}
	user.Username = username
	if err := row.Scan(&user.ID, &user.Balance); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *mariaDBStore) GetUsers() ([]*models.User, error) {
	rows, err := s.db.Query("SELECT id, username, balance FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.ID, &user.Username, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *mariaDBStore) SearchUsers(searchTerm string) ([]*models.User, error) {
	searchTermFormatted := fmt.Sprintf("%%%s%%", searchTerm)
	query := `
		SELECT id, username, balance
		FROM users
		WHERE username LIKE ?
    `
	rows, err := s.db.Query(query, searchTermFormatted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.ID, &user.Username, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *mariaDBStore) Insert(u *models.User) error {
	result, err := s.db.Exec("INSERT INTO users (id, username, balance) VALUES (?, ?, ?)", u.ID, u.Username, u.Balance)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("error inserting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows affected after insert: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("insert failed: no rows affected")
	}

	return nil
}

func (s *mariaDBStore) Update(u *models.User) error {
	result, err := s.db.Exec("UPDATE users SET username=? WHERE id=?", u.Username, u.ID)
	if err != nil {
		return fmt.Errorf("error executing user update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows affected after update: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotExists
	}

	return nil
}

func (s *mariaDBStore) Remove(id int) error {
	result, err := s.db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error executing user delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows affected after update: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotExists
	}

	return nil
}

func (s *mariaDBStore) Exists(id int) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE id=?", id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return count > 0, nil
}
