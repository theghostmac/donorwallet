package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the donor wallet system.
type User struct {
	UserID uuid.UUID `json:"user_id"`
	Username string `json:"username"`
	PasswordHash string `json:"password_hash"` // not supposed to expose in JSON response.
	TxnPIN string `json:"txn_pin"` // not suppose to expose in JSON response
	Email string `json:"email"`
	WalletBalance float64 `json:"wallet_balance,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Wallet represents a user's wallet.
type Wallet struct {
	WalletID uint `json:"wallet_id"`
	UserID uuid.UUID `json:"user_id"`
	Balance float64 `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	TxnID uint `json:"txn_id"`
	UserID uuid.UUID `json:"user_id"`
	Amount float64 `json:"amount"`
	Type string `json:"type"` // "deposit", "withdraw", "donation"
	CreatedAt time.Time `json:"created_at"`
}

type Donation struct {
	DonationID uint `gorm:"primary_key" json:"donation_id"`
	DonorID uuid.UUID `json:"donor_id"`
	BeneficiaryID  uuid.UUID `json:"beneficiary_id"`
	Amount float64 `json:"amount"`
	Message string `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}