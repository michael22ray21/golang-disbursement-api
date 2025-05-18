package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// handleIndex handles the index route
func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Status:  "available",
		Message: "Paper.id Disbursement API is running",
	}

	app.writeJSON(w, http.StatusOK, response)
}

// handleDisbursement handles the disbursement request
func (app *App) handleDisbursement(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req DisbursementRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.handleError(w, ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	// Validate request
	err = app.Validator.Struct(req)
	if err != nil {
		app.handleError(w, fmt.Errorf("%w: %v", ErrInvalidRequest, err), http.StatusBadRequest)
		return
	}

	// Begin transaction
	tx, err := app.DB.Begin()
	if err != nil {
		app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check if user exists
	var user User
	err = tx.QueryRow("SELECT id, name, email FROM users WHERE id = ?", req.UserID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			app.handleError(w, ErrUserNotFound, http.StatusNotFound)
		} else {
			app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		}
		return
	}

	// Check if bank account exists and belongs to the user
	var bankAccount BankAccount
	err = tx.QueryRow(
		"SELECT id, user_id, bank_name, account_number, account_name FROM bank_accounts WHERE id = ? AND user_id = ?",
		req.BankAccountID, req.UserID,
	).Scan(&bankAccount.ID, &bankAccount.UserID, &bankAccount.BankName, &bankAccount.AccountNumber, &bankAccount.AccountName)
	if err != nil {
		if err == sql.ErrNoRows {
			app.handleError(w, ErrBankAccountNotFound, http.StatusNotFound)
		} else {
			app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		}
		return
	}

	// Check if wallet exists and has sufficient balance
	var wallet Wallet
	err = tx.QueryRow("SELECT id, user_id, balance, currency FROM wallets WHERE user_id = ?", req.UserID).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			app.handleError(w, ErrWalletNotFound, http.StatusNotFound)
		} else {
			app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		}
		return
	}

	// Check if balance is sufficient
	if wallet.Balance < req.Amount {
		app.handleError(w, ErrInsufficientBalance, http.StatusBadRequest)
		return
	}

	// Generate a unique reference number for the disbursement
	reference := fmt.Sprintf("DSB%d%d", time.Now().Unix(), req.UserID)

	// Create disbursement record
	now := time.Now()
	result, err := tx.Exec(
		`INSERT INTO disbursements
		(user_id, wallet_id, bank_account_id, amount, status, reference, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.UserID, wallet.ID, req.BankAccountID, req.Amount, "pending", reference, req.Description, now, now,
	)
	if err != nil {
		app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	disbursementID, err := result.LastInsertId()
	if err != nil {
		app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// Update wallet balance
	_, err = tx.Exec("UPDATE wallets SET balance = balance - ?, updated_at = ? WHERE id = ?", req.Amount, now, wallet.ID)
	if err != nil {
		app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// In actual disbursement scenario, we would call an external disbursement service or third-party service here
	// For this test, we'll simulate a successful disbursement
	completedAt := time.Now()
	_, err = tx.Exec(
		"UPDATE disbursements SET status = ?, completed_at = ?, updated_at = ? WHERE id = ?",
		"completed", completedAt, completedAt, disbursementID,
	)
	if err != nil {
		app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// Get the disbursement to build response
	disbursement, err := app.getDisbursementByID(disbursementID)
	if err != nil {
		app.handleError(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// Prepare response data
	response := DisbursementResponse{
		Success:      true,
		Message:      "Disbursement processed successfully",
		Disbursement: disbursement,
	}

	// Send response
	app.writeJSON(w, http.StatusOK, response)
}

// getDisbursementByID retrieves a disbursement by ID
func (app *App) getDisbursementByID(id int64) (*Disbursement, error) {
	var disbursement Disbursement
	var completedAt sql.NullTime

	err := app.DB.QueryRow(`
		SELECT id, user_id, wallet_id, bank_account_id, amount, status, reference, created_at, updated_at, completed_at
		FROM disbursements
		WHERE id = ?
	`, id).Scan(
		&disbursement.ID,
		&disbursement.UserID,
		&disbursement.WalletID,
		&disbursement.BankAccountID,
		&disbursement.Amount,
		&disbursement.Status,
		&disbursement.Reference,
		&disbursement.CreatedAt,
		&disbursement.UpdatedAt,
		&completedAt,
	)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		disbursement.CompletedAt = &completedAt.Time
	}

	return &disbursement, nil
}

// writeJSON writes JSON response
func (app *App) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// handleError handles error responses
func (app *App) handleError(w http.ResponseWriter, err error, status int) {
	log.Printf("Error: %v", err)
	response := DisbursementResponse{
		Success: false,
		Message: "An error occurred",
		Error:   err.Error(),
	}
	app.writeJSON(w, status, response)
}
