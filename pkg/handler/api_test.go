package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/lattots/piikittaja/pkg/handler"
	"github.com/lattots/piikittaja/pkg/models"
	"github.com/lattots/piikittaja/pkg/transaction"
	userstore "github.com/lattots/piikittaja/pkg/user_store"
)

func getTestRouter(testHandler *handler.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", testHandler.GetUsers)
	mux.HandleFunc("GET /users/{userId}", testHandler.GetUserByID)
	mux.HandleFunc("POST /users/{userId}/transactions", testHandler.NewTransaction)

	return mux
}

func TestGetUserByID(t *testing.T) {
	usrStore, _ := userstore.NewMockStore()
	traStore := transaction.NewMockStore()
	traHandler := transaction.NewTransactionHandler(traStore)

	// Populate user store with some test data

	testUser := &models.User{
		ID:       123,
		Username: "Foo",
	}
	_ = usrStore.Insert(testUser)

	testHandler := handler.NewTestHandler(usrStore, traHandler)
	router := getTestRouter(testHandler)

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{"Get Existing User", "GET", "/users/123", http.StatusOK},
		{"Get User That Doesn't Exist", "GET", "/users/456", http.StatusNotFound},
		{"Invalid Method", "POST", "/users/123", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, rr.Code)
			}

			// Early return for tests that should fail
			if rr.Code != http.StatusOK {
				return
			}

			var respUser models.UserResponse
			err := json.NewDecoder(rr.Body).Decode(&respUser)
			if err != nil {
				t.Errorf("%s: unexpected error when decoding response: %s", tt.name, err)
			}

			if respUser.Username != testUser.Username {
				t.Errorf("%s: wrong username for test user: want %s got %s", tt.name, testUser.Username, respUser.Username)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	dbUrl := os.Getenv("DATABASE_APP")
	if dbUrl == "" {
		t.Fatal("DATABASE_APP not set")
	}

	// Test MariaDB store should already be populated
	usrStore, _ := userstore.NewMariaDBStore(dbUrl)
	defer usrStore.Close()
	traStore, _ := transaction.NewMariaDBStore(dbUrl)
	defer traStore.Close()
	traHandler := transaction.NewTransactionHandler(traStore)

	testHandler := handler.NewTestHandler(usrStore, traHandler)
	router := getTestRouter(testHandler)

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{"Get Users", "GET", "/users", http.StatusOK},
		{"Invalid Method", "POST", "/users", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, rr.Code)
			}

			// Early return for tests that should fail
			if rr.Code != http.StatusOK {
				return
			}

			var respUsers []models.UserResponse
			err := json.NewDecoder(rr.Body).Decode(&respUsers)
			if err != nil {
				t.Errorf("%s: unexpected error when decoding response: %s", tt.name, err)
			}

			if len(respUsers) == 0 {
				t.Errorf("%s: no users in response", tt.name)
			}
		})
	}
}

func TestNewTransaction(t *testing.T) {
	dbUrl := os.Getenv("DATABASE_APP")
	if dbUrl == "" {
		t.Fatal("DATABASE_APP not set")
	}

	// Test MariaDB store should already be populated
	usrStore, _ := userstore.NewMariaDBStore(dbUrl)
	defer usrStore.Close()
	traStore, _ := transaction.NewMariaDBStore(dbUrl)
	defer traStore.Close()
	traHandler := transaction.NewTransactionHandler(traStore)

	testHandler := handler.NewTestHandler(usrStore, traHandler)
	router := getTestRouter(testHandler)

	originalUsers, _ := usrStore.GetUsers()
	if len(originalUsers) == 0 {
		t.Fatal("failed to get users")
	}
	testUser := originalUsers[0]

	const testAmount int = 1000
	reqContent := fmt.Sprintf(`{"type":"deposit", "amount":%d}`, testAmount)

	req := httptest.NewRequest("POST", fmt.Sprintf("/users/%d/transactions", testUser.ID), strings.NewReader(reqContent))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	var respUser models.UserResponse
	err := json.NewDecoder(rr.Body).Decode(&respUser)
	if err != nil {
		t.Errorf("unexpected error when decoding response: %s", err)
	}

	expectedBalance := testUser.Balance + testAmount
	if respUser.Balance != expectedBalance {
		t.Errorf("wrong balance in response: want %d got %d", expectedBalance, respUser.Balance)
	}

	reqContent = fmt.Sprintf(`{"type":"withdraw", "amount":%d}`, testAmount)

	req = httptest.NewRequest("POST", fmt.Sprintf("/users/%d/transactions", testUser.ID), strings.NewReader(reqContent))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	err = json.NewDecoder(rr.Body).Decode(&respUser)
	if err != nil {
		t.Errorf("unexpected error when decoding response: %s", err)
	}

	if respUser.Balance != testUser.Balance {
		t.Errorf("wrong balance in response: want %d got %d", testUser.Balance, respUser.Balance)
	}
}
