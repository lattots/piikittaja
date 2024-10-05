package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/lattots/piikittaja/pkg/auth"
	"github.com/lattots/piikittaja/pkg/env"
)

type Handler struct {
	DB   *sql.DB
	tmpl *template.Template
	Auth *auth.Service
}

const (
	cookieStoreMaxAge = 86400 * 30
	cookieStoreIsProd = true
)

func NewHandler() (*Handler, error) {
	// Database URL is read from environment variables.
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		return nil, errors.New("error getting database URL from environment variables")
	}

	// Database handle is created.
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, errors.New("error connecting to the database")
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	cookieStoreKey := os.Getenv("COOKIE_STORE_SECRET")

	sessionOptions := auth.SessionOptions{
		CookiesKey: cookieStoreKey,
		MaxAge:     cookieStoreMaxAge,
		Secure:     cookieStoreIsProd,
	}

	store := auth.NewCookieStore(sessionOptions)

	authService := auth.NewService(store, db)

	root, err := env.GetProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("error getting project root folder: %w", err)
	}

	// HTML template file is parsed
	tmpl, err := template.ParseFiles(filepath.Join(root, "assets/web_app/html/template.html"))
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %s", err)
	}

	return &Handler{DB: db, tmpl: tmpl, Auth: authService}, nil
}
