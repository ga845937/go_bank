package walletModule

import (
	"github.com/gin-gonic/gin"

	"go_bank/internal/middleware"
)

type WalletRoute struct {
	walletController *WalletController
}

func NewWalletRoute(walletController *WalletController) *WalletRoute {
	return &WalletRoute{
		walletController,
	}
}

func (walletRoute *WalletRoute) RegisterRoute(walletRouteGroup *gin.RouterGroup) {
	walletRouteGroup.GET("/", middleware.Validator(walletRoute.walletController.ReadWallet))
	walletRouteGroup.POST("/operation", middleware.Validator(walletRoute.walletController.Operation))
}
