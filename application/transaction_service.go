package application

import (
	"EWallet/domain"
	"errors"
	"fmt"
	"time"
)

type TransactionService struct {
	transactionRepository domain.TransactionRepositoryInterface
	walletService         domain.WalletServiceInterface
}

func NewTransactionService(transactionRepository domain.TransactionRepositoryInterface, walletService domain.WalletServiceInterface) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		walletService:         walletService,
	}
}

func (s *TransactionService) CreateTransaction(fromWalletID, toWalletID string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	fromWallet, err := s.walletService.GetWalletByID(fromWalletID)
	if err != nil {
		return err
	}
	if fromWallet.Balance < amount {
		return errors.New("insufficient balance")
	}
	if fromWalletID == toWalletID {
		return errors.New("impossible to self send")
	}

	toWallet, err := s.walletService.GetWalletByID(toWalletID)
	if err != nil {
		return err
	}

	transaction := &domain.Transaction{
		Time:   time.Now(),
		From:   fromWalletID,
		To:     toWalletID,
		Amount: amount,
	}

	err = s.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	err = s.walletService.UpdateWalletBalance(fromWalletID, fromWallet.Balance-amount)
	if err != nil {
		return err
	}

	err = s.walletService.UpdateWalletBalance(toWalletID, toWallet.Balance+amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) GetTransactionsByWalletID(walletID string) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepository.GetTransactionsByWalletID(walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by wallet ID: %w", err)
	}
	return transactions, nil
}
