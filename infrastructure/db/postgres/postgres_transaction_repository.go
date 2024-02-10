package postgres

import (
	"EWallet/domain"
	"database/sql"
	"fmt"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetTransactionsByWalletID(walletID string) ([]*domain.Transaction, error) {
	query := `
		SELECT time, from_wallet_id, to_wallet_id, amount
		FROM transactions
		WHERE from_wallet_id = $1 OR to_wallet_id = $1
	`
	rows, err := r.db.Query(query, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by wallet ID: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(&transaction.Time, &transaction.From, &transaction.To, &transaction.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) CreateTransaction(transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (time, from_wallet_id, to_wallet_id, amount)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, transaction.Time, transaction.From, transaction.To, transaction.Amount)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}
