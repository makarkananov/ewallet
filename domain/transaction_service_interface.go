package domain

type TransactionServiceInterface interface {
	CreateTransaction(fromWalletID, toWalletID string, amount float64) error
	GetTransactionsByWalletID(walletID string) ([]*Transaction, error)
}
