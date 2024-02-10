package domain

type WalletRepositoryInterface interface {
	GetWalletByID(walletID string) (*Wallet, error)
	CreateWallet(wallet *Wallet) error
	UpdateWalletBalance(walletID string, newBalance float64) error
}
