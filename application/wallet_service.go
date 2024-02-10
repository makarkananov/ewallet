package application

import (
	"EWallet/domain"
	"fmt"
	"github.com/google/uuid"
)

type WalletService struct {
	walletRepository domain.WalletRepositoryInterface
}

func NewWalletService(walletRepository domain.WalletRepositoryInterface) *WalletService {
	return &WalletService{
		walletRepository: walletRepository,
	}
}

func (s *WalletService) CreateWallet() (*domain.Wallet, error) {
	walletID := uuid.New().String()

	wallet := &domain.Wallet{
		ID:      walletID,
		Balance: 100.0,
	}

	err := s.walletRepository.CreateWallet(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}
	return wallet, nil
}

func (s *WalletService) GetWalletByID(walletID string) (*domain.Wallet, error) {
	wallet, err := s.walletRepository.GetWalletByID(walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet by ID: %w", err)
	}
	return wallet, nil
}

func (s *WalletService) UpdateWalletBalance(id string, balance float64) error {
	err := s.walletRepository.UpdateWalletBalance(id, balance)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}
	return nil
}
