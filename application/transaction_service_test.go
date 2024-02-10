package application

import (
	"EWallet/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetWalletByID(walletID string) (*domain.Wallet, error) {
	args := m.Called(walletID)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletService) UpdateWalletBalance(id string, balance float64) error {
	args := m.Called(id, balance)
	return args.Error(0)
}

func (m *MockWalletService) CreateWallet() (*domain.Wallet, error) {
	args := m.Called()
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetTransactionsByWalletID(walletID string) ([]*domain.Transaction, error) {
	args := m.Called(walletID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) CreateTransaction(transaction *domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func TestTransactionService_GetTransactionsByWalletID(t *testing.T) {
	// Arrange
	mockRepo := &MockTransactionRepository{}
	service := NewTransactionService(mockRepo, nil)
	walletID := "testWalletID"
	mockTransactions := []*domain.Transaction{{Time: time.Now(), From: "from", To: "to", Amount: 50.0}}
	mockRepo.On("GetTransactionsByWalletID", walletID).Return(mockTransactions, nil)

	// Act
	transactions, err := service.GetTransactionsByWalletID(walletID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, mockTransactions, transactions)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_CreateTransaction(t *testing.T) {
	// Arrange
	mockRepo := &MockTransactionRepository{}
	mockWalletService := &MockWalletService{}
	service := NewTransactionService(mockRepo, mockWalletService)
	fromWalletID := "fromWallet"
	toWalletID := "toWallet"
	amount := 50.0
	mockWalletService.On("GetWalletByID", fromWalletID).Return(&domain.Wallet{Balance: 100.0}, nil)
	mockWalletService.On("GetWalletByID", toWalletID).Return(&domain.Wallet{Balance: 80.0}, nil)
	mockRepo.On("CreateTransaction", mock.Anything).Return(nil)
	mockWalletService.On("UpdateWalletBalance", mock.Anything, mock.Anything).Return(nil).Times(2)

	// Act
	err := service.CreateTransaction(fromWalletID, toWalletID, amount)

	// Assert
	assert.NoError(t, err)
	mockWalletService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
