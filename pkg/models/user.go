package models

type User struct {
	ID       int
	Username string
	Balance  int // Balance in cents
	IsAdmin  bool
}

func NewUser(id int, username string) *User {
	return &User{ID: id, Username: username, Balance: 0}
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}

type TransactionRequest struct {
	UserID int    `json:"userId"`
	Issuer string `json:"issuer"`
	Type   string `json:"type"`   // Transaction type: withdraw / deposit
	Amount int    `json:"amount"` // Transaction amount in cents
}
