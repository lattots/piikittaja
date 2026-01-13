package handler

import (
	"errors"
	"fmt"
	"os"

	"github.com/lattots/piikittaja/pkg/auth"
	"github.com/lattots/piikittaja/pkg/transaction"
	userstore "github.com/lattots/piikittaja/pkg/user_store"
)

type Handler struct {
	hostUrl    string
	traHandler transaction.TransactionHandler
	usrStore   userstore.UserStore
	Auth       *auth.Service
}

const (
	cookieStoreMaxAge = 86400 * 30
	cookieStoreIsProd = true
)

func NewHandler(hostUrl string) (*Handler, error) {
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

	h := &Handler{
		hostUrl:    hostUrl,
		traHandler: traHandler,
		usrStore:   usrStore,
		Auth:       authService,
	}

	return h, nil
}

func NewTestHandler(usrStore userstore.UserStore, traHandler transaction.TransactionHandler) *Handler {
	return &Handler{
		usrStore:   usrStore,
		traHandler: traHandler,
	}
}
