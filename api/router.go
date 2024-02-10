package api

import "github.com/gin-gonic/gin"

func SetupRouter(api *API) *gin.Engine {
	router := gin.Default()

	router.POST("/api/v1/wallet", api.CreateWallet)
	router.GET("/api/v1/wallet/:walletId", api.GetWalletByID)
	router.POST("/api/v1/wallet/:walletId/send", api.SendMoney)
	router.GET("/api/v1/wallet/:walletId/history", api.GetTransactionHistory)

	return router
}
