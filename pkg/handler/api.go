package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/lattots/piikittaja/pkg/models"
	"github.com/lattots/piikittaja/pkg/transaction"
	userstore "github.com/lattots/piikittaja/pkg/user_store"
)

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		http.Error(w, "No userId provided in request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "userId is not an integer", http.StatusBadRequest)
		return
	}

	usr, err := h.usrStore.GetByID(id)
	if errors.Is(err, userstore.ErrUserNotExists) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	usrResp := usr.ToResponse()

	writeJSONResponse(w, http.StatusOK, usrResp)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("searchTerm")
	var usrs []*models.User
	var err error
	if searchTerm == "" {
		usrs, err = h.usrStore.GetUsers()
		if err != nil {
			log.Printf("error getting users from store: %s\n", err)
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}
	} else {
		usrs, err = h.usrStore.SearchUsers(searchTerm)
		if err != nil {
			log.Printf("error searching users in store: %s\n", err)
			http.Error(w, "Failed to search for users", http.StatusInternalServerError)
			return
		}
	}

	respContent := make([]models.UserResponse, len(usrs))
	for i := range usrs {
		respContent[i] = usrs[i].ToResponse()
	}

	writeJSONResponse(w, http.StatusOK, respContent)
}

func (h *Handler) NewTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		http.Error(w, "No userId provided in request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "userId is not an integer", http.StatusBadRequest)
		return
	}

	usr, err := h.usrStore.GetByID(id)
	if errors.Is(err, userstore.ErrUserNotExists) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content must be transaction request JSON", http.StatusUnsupportedMediaType)
		return
	}

	var req models.Transaction
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("error decoding transaction request: %s\n", err)
		http.Error(w, "Failed to decode request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	switch req.Type {
	case "deposit":
		_, err = h.traHandler.Deposit(usr, req.Amount)
	case "withdraw":
		_, err = h.traHandler.Withdraw(usr, req.Amount)
	}
	if errors.Is(err, transaction.ErrNotEnoughBalance) {
		http.Error(w, "User doesn't have enough balance", http.StatusPaymentRequired)
		return
	}
	if errors.Is(err, transaction.ErrInvalidAmount) {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("error making transaction: %s\n", err)
		http.Error(w, "Failed to execute transaction", http.StatusInternalServerError)
		return
	}

	usrResp := usr.ToResponse()

	writeJSONResponse(w, http.StatusCreated, usrResp)
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2)
	endDateStr := r.URL.Query().Get("endDate")
	if endDateStr == "" {
		http.Error(w, "no end date in query, please provide \"endDate\"", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		http.Error(w, "invalid end date format", http.StatusBadRequest)
		return
	}

	windowStr := r.URL.Query().Get("window")
	if windowStr == "" {
		http.Error(w, "no transaction time window in query, please provide \"window\"", http.StatusBadRequest)
		return
	}
	windowInt, err := strconv.Atoi(windowStr)
	if err != nil {
		http.Error(w, "invalid window in query", http.StatusBadRequest)
		return
	}
	window := time.Duration(windowInt) * 24 * time.Hour // Window in request represents the number of days

	traType := r.URL.Query().Get("type")

	transactions, err := h.traHandler.GetTransactions(endDate, window, traType)
	if err != nil {
		log.Printf("error getting transactions: %s\n", err)
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, transactions)
}

const defaultTransactionQuantity = 3

func (h *Handler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		http.Error(w, "No userId provided in request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "userId is not an integer", http.StatusBadRequest)
		return
	}

	usr, err := h.usrStore.GetByID(id)
	if errors.Is(err, userstore.ErrUserNotExists) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	quantityStr := r.URL.Query().Get("quantity")
	quantity := defaultTransactionQuantity
	if quantityStr != "" {
		quantityInt, err := strconv.Atoi(quantityStr)
		if err == nil {
			quantity = quantityInt
		}
	}

	transactions, err := h.traHandler.GetUserTransactions(usr, quantity)
	if err != nil {
		log.Printf("error getting transactions for user %d: %s\n", usr.ID, err)
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, transactions)
}

// writeJSONResponse writes JSON encodable data to response writer with the provided status code
func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error writing JSON response: %s\n", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
