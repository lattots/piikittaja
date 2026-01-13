package models

type TransactionRequest struct {
	Issuer string `json:"issuer"`
	Type   string `json:"type"`   // Transaction type: withdraw / deposit
	Amount int    `json:"amount"` // Transaction amount in cents
}
