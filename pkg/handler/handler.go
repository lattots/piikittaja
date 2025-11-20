package handler

import (
	"errors"
	"fmt"
	"html/template"
	"os"

	"github.com/lattots/piikittaja/pkg/auth"
	"github.com/lattots/piikittaja/pkg/transaction"
	userstore "github.com/lattots/piikittaja/pkg/user_store"
)

type Handler struct {
	traHandler transaction.TransactionHandler
	usrStore   userstore.UserStore
	tmpl       *template.Template
	Auth       *auth.Service
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

	usrStore, err := userstore.NewMariaDBStore(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating user store: %w", err)
	}

	traStore, err := transaction.NewMariaDBStore(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction store: %w", err)
	}
	traHandler := transaction.NewTransactionHandler(traStore)

	cookieStoreKey := os.Getenv("COOKIE_STORE_SECRET")

	sessionOptions := auth.SessionOptions{
		CookiesKey: cookieStoreKey,
		MaxAge:     cookieStoreMaxAge,
		Secure:     cookieStoreIsProd,
	}

	store := auth.NewCookieStore(sessionOptions)

	authService, err := auth.NewService(store, dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating auth service: %w", err)
	}

	// HTML template file is parsed
	tmpl, err := template.ParseFiles("./assets/web_app/html/template.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %s", err)
	}

	h := &Handler{
		traHandler: traHandler,
		usrStore:   usrStore,
		tmpl:       tmpl,
		Auth:       authService,
	}

	return h, nil
}
