package postgres

import (
	"EWallet/domain"
	"database/sql"
	"fmt"
)

type WalletRepository struct {
	db *sql.DB
}

func NewPostgresWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetWalletByID(walletID string) (*domain.Wallet, error) {
	query := "SELECT id, balance FROM wallets WHERE id = $1"
	row := r.db.QueryRow(query, walletID)

	var wallet domain.Wallet
	err := row.Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet by ID: %w", err)
	}

	return &wallet, nil
}

func (r *WalletRepository) CreateWallet(wallet *domain.Wallet) error {
	query := "INSERT INTO wallets (id, balance) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(query, wallet.ID, wallet.Balance).Scan(&wallet.ID)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (r *WalletRepository) UpdateWalletBalance(walletID string, newBalance float64) error {
	query := "UPDATE wallets SET balance = $1 WHERE id = $2"
	_, err := r.db.Exec(query, newBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	return nil
}
