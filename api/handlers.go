package api

import (
	"EWallet/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type API struct {
	walletService      domain.WalletServiceInterface
	transactionService domain.TransactionServiceInterface
}

func NewAPI(walletService domain.WalletServiceInterface, transactionService domain.TransactionServiceInterface) *API {
	return &API{
		walletService:      walletService,
		transactionService: transactionService,
	}
}

func (api *API) CreateWallet(c *gin.Context) {
	wallet, err := api.walletService.CreateWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (api *API) GetWalletByID(c *gin.Context) {
	walletID := c.Param("walletId")
	wallet, err := api.walletService.GetWalletByID(walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func (api *API) SendMoney(c *gin.Context) {
	fromWalletID := c.Param("walletId")
	var request struct {
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.transactionService.CreateTransaction(fromWalletID, request.To, request.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction successful"})
}

func (api *API) GetTransactionHistory(c *gin.Context) {
	walletID := c.Param("walletId")
	transactions, err := api.transactionService.GetTransactionsByWalletID(walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
