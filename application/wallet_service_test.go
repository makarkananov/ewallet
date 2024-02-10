package application

import (
	"EWallet/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetWalletByID(walletID string) (*domain.Wallet, error) {
	args := m.Called(walletID)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletRepository) CreateWallet(wallet *domain.Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) UpdateWalletBalance(walletID string, newBalance float64) error {
	args := m.Called(walletID, newBalance)
	return args.Error(0)
}

func TestWalletService_CreateWallet(t *testing.T) {
	// Arrange
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)
	mockRepo.On("CreateWallet", mock.Anything).Return(nil)

	// Act
	wallet, err := service.CreateWallet()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, wallet.Balance, 100.0)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_GetWalletByID(t *testing.T) {
	// Arrange
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)
	walletID := "testID"
	mockWallet := &domain.Wallet{ID: walletID, Balance: 100.0}
	mockRepo.On("GetWalletByID", walletID).Return(mockWallet, nil)

	// Act
	wallet, err := service.GetWalletByID(walletID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, mockWallet, wallet)
	mockRepo.AssertExpectations(t)
}

func TestWalletService_UpdateWalletBalance(t *testing.T) {
	// Arrange
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)
	walletID := "testID"
	newBalance := 150.0
	mockRepo.On("UpdateWalletBalance", walletID, newBalance).Return(nil)

	// Act
	err := service.UpdateWalletBalance(walletID, newBalance)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
