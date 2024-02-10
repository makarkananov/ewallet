package domain

type WalletServiceInterface interface {
	CreateWallet() (*Wallet, error)
	GetWalletByID(walletID string) (*Wallet, error)
	UpdateWalletBalance(id string, balance float64) error
}
