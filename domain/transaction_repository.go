package domain

type TransactionRepositoryInterface interface {
	GetTransactionsByWalletID(walletID string) ([]*Transaction, error)
	CreateTransaction(transaction *Transaction) error
}
