package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// initializeDB sets up the database connection and schema
func initializeDB(dsn string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Ensure tables exist
	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// createTables creates the necessary database tables if they don't exist
func createTables(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create bank_accounts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bank_accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			bank_name TEXT NOT NULL,
			account_number TEXT NOT NULL,
			account_name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		return err
	}

	// Create wallets table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS wallets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			balance INTEGER NOT NULL DEFAULT 0,
			currency TEXT NOT NULL DEFAULT 'IDR',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		return err
	}

	// Create disbursements table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS disbursements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			wallet_id INTEGER NOT NULL,
			bank_account_id INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			status TEXT NOT NULL,
			reference TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (wallet_id) REFERENCES wallets (id),
			FOREIGN KEY (bank_account_id) REFERENCES bank_accounts (id)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// seedTestData populates the database with test data
func seedTestData(app *App) error {
	// Check if users already exist
	var count int
	err := app.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	// Skip seeding if data already exists
	if count > 0 {
		return nil
	}

	// Begin transaction
	tx, err := app.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert test users
	users := []struct {
		Name  string
		Email string
	}{
		{"John Doe", "john@example.com"},
		{"Jane Smith", "jane@example.com"},
		{"Mike Seed", "mike@example.com"},
	}

	for _, user := range users {
		_, err = tx.Exec(
			"INSERT INTO users (name, email) VALUES (?, ?)",
			user.Name, user.Email,
		)
		if err != nil {
			return err
		}
	}

	// Insert bank accounts
	bankAccounts := []struct {
		UserID        int64
		BankName      string
		AccountNumber string
		AccountName   string
	}{
		{1, "BCA", "1234567890", "John Doe"},
		{2, "Mandiri", "0987654321", "Jane Smith"},
		{3, "BNI", "1122334455", "Mike Seed"},
	}

	for _, bankAccount := range bankAccounts {
		_, err = tx.Exec(
			"INSERT INTO bank_accounts (user_id, bank_name, account_number, account_name) VALUES (?, ?, ?, ?)",
			bankAccount.UserID, bankAccount.BankName, bankAccount.AccountNumber, bankAccount.AccountName,
		)
		if err != nil {
			return err
		}
	}

	// Insert wallets with initial balances
	wallets := []struct {
		UserID   int64
		Balance  int64
		Currency string
	}{
		{1, 5000000, "IDR"}, // Rp 5,000,000
		{2, 2500000, "IDR"}, // Rp 2,500,000
		{3, 7500000, "IDR"}, // Rp 7,500,000
	}

	for _, wallet := range wallets {
		_, err = tx.Exec(
			"INSERT INTO wallets (user_id, balance, currency) VALUES (?, ?, ?)",
			wallet.UserID, wallet.Balance, wallet.Currency,
		)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}
