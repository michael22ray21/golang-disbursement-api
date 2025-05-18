package main

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	BankAccountID *int64    `json:"bank_account_id,omitempty"`
}

// Wallet represents a user's wallet
type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   int64     `json:"balance"` // Balance is in smallest currency unit in cases of cents and dollars
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BankAccount represents a user's bank account for disbursement
type BankAccount struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	BankName      string    `json:"bank_name"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Disbursement represents a disbursement transaction
type Disbursement struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	WalletID      int64      `json:"wallet_id"`
	BankAccountID int64      `json:"bank_account_id"`
	Amount        int64      `json:"amount"` // Amount in smallest currency unit
	Status        string     `json:"status"` // pending, completed, failed
	Reference     string     `json:"reference"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

// DisbursementRequest represents the payload for a disbursement request
type DisbursementRequest struct {
	UserID        int64  `json:"user_id" validate:"required"`
	BankAccountID int64  `json:"bank_account_id" validate:"required"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Description   string `json:"description" validate:"omitempty,max=255"`
}

// DisbursementResponse represents the response for a disbursement request
type DisbursementResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	Disbursement *Disbursement `json:"disbursement,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// StatusResponse represents a generic status response
type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
