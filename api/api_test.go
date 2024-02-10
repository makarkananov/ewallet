package api

import (
	"EWallet/domain"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet() (*domain.Wallet, error) {
	args := m.Called()
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletService) GetWalletByID(walletID string) (*domain.Wallet, error) {
	args := m.Called(walletID)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletService) UpdateWalletBalance(id string, balance float64) error {
	args := m.Called(id, balance)
	return args.Error(0)
}

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(fromWalletID, toWalletID string, amount float64) error {
	args := m.Called(fromWalletID, toWalletID, amount)
	return args.Error(0)
}

func (m *MockTransactionService) GetTransactionsByWalletID(walletID string) ([]*domain.Transaction, error) {
	args := m.Called(walletID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func TestCreateWallet(t *testing.T) {
	// Arrange
	mockWalletService := &MockWalletService{}
	mockTransactionService := &MockTransactionService{}
	mockWalletService.On("CreateWallet").Return(&domain.Wallet{}, nil)
	newAPI := NewAPI(mockWalletService, mockTransactionService)
	router := SetupRouter(newAPI)

	// Act
	req, err := http.NewRequest("POST", "/api/v1/wallet", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	mockWalletService.AssertExpectations(t)
}

func TestGetWalletByID(t *testing.T) {
	// Arrange
	mockWalletService := &MockWalletService{}
	mockTransactionService := &MockTransactionService{}
	mockWalletService.On("GetWalletByID", "testID").Return(&domain.Wallet{}, nil)
	api := NewAPI(mockWalletService, mockTransactionService)
	router := SetupRouter(api)

	// Act
	req, err := http.NewRequest("GET", "/api/v1/wallet/testID", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	mockWalletService.AssertExpectations(t)
}

func TestSendMoney(t *testing.T) {
	// Arrange
	mockWalletService := &MockWalletService{}
	mockTransactionService := &MockTransactionService{}
	mockTransactionService.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	api := NewAPI(mockWalletService, mockTransactionService)
	router := SetupRouter(api)

	// Act
	requestBody := []byte(`{"to": "recipientID", "amount": 50.0}`)
	req, err := http.NewRequest("POST", "/api/v1/wallet/senderID/send", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	mockTransactionService.AssertExpectations(t)
}

func TestGetTransactionHistory(t *testing.T) {
	// Arrange
	mockWalletService := &MockWalletService{}
	mockTransactionService := &MockTransactionService{}
	mockTransactionService.On("GetTransactionsByWalletID", "testID").Return([]*domain.Transaction{}, nil)
	api := NewAPI(mockWalletService, mockTransactionService)
	router := SetupRouter(api)

	// Act
	req, err := http.NewRequest("GET", "/api/v1/wallet/testID/history", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	mockTransactionService.AssertExpectations(t)
}

func TestCreateWallet_Error(t *testing.T) {
	// Arrange
	mockWalletService := &MockWalletService{}
	mockTransactionService := &MockTransactionService{}
	expectedWallet := &domain.Wallet{}
	expectedError := errors.New("error creating wallet")

	mockWalletService.On("CreateWallet").Return(expectedWallet, expectedError)

	newAPI := NewAPI(mockWalletService, mockTransactionService)
	router := SetupRouter(newAPI)

	// Act
	req, err := http.NewRequest("POST", "/api/v1/wallet", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), expectedError.Error())

	mockWalletService.AssertExpectations(t)
}

func TestGetWalletByID_NotFound(t *testing.T) {
	// Arrange
	mockWalletService := &MockWalletService{}
	mockTransactionService := &MockTransactionService{}
	expectedError := errors.New("wallet not found")
	expectedWallet := &domain.Wallet{}

	mockWalletService.On("GetWalletByID", "nonExistentID").Return(expectedWallet, expectedError)
	api := NewAPI(mockWalletService, mockTransactionService)
	router := SetupRouter(api)

	// Act
	req, err := http.NewRequest("GET", "/api/v1/wallet/nonExistentID", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), expectedError.Error())
	mockWalletService.AssertExpectations(t)
}

func TestSendMoney_InvalidRequest(t *testing.T) {
	// Arrange
	mockWalletService := &MockWalletService{}
	mockTransactionService := &MockTransactionService{}
	api := NewAPI(mockWalletService, mockTransactionService)
	router := SetupRouter(api)

	// Act
	req, err := http.NewRequest("POST", "/api/v1/wallet/senderID/send", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid request")
	mockTransactionService.AssertNotCalled(t, "CreateTransaction")
}

func TestSendMoney_InsufficientFunds(t *testing.T) {
	// Arrange
	mockWalletService := &MockWalletService{}
	mockTransactionService := &MockTransactionService{}
	expectedError := errors.New("insufficient funds")
	mockTransactionService.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything).Return(expectedError)
	api := NewAPI(mockWalletService, mockTransactionService)
	router := SetupRouter(api)

	// Act
	requestBody := []byte(`{"to": "recipientID", "amount": 50.0}`)
	req, err := http.NewRequest("POST", "/api/v1/wallet/senderID/send", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), expectedError.Error())
	mockTransactionService.AssertExpectations(t)
}
