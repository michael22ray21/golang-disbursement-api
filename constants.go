package main

import "errors"

// Common errors
var (
	ErrInvalidJSON         = errors.New("invalid JSON payload")
	ErrInvalidRequest      = errors.New("invalid request")
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrUserNotFound        = errors.New("user not found")
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrBankAccountNotFound = errors.New("bank account not found")
	ErrInternalServer      = errors.New("internal server error")
)
