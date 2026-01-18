package models

import "time"

type Transaction struct {
	Issuer   string    `json:"issuer"`
	IssuedAt time.Time `json:"issuedAt"` // Transaction timestamp as RFC3339
	Type     string    `json:"type"`     // Transaction type: withdraw / deposit
	Amount   int       `json:"amount"`   // Transaction amount in cents
}
