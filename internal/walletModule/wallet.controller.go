package walletModule

import (
	"github.com/gin-gonic/gin"

	middleware "go_bank/internal/middleware"
	walletEntity "go_bank/internal/walletModule/entity"
	walletProvider "go_bank/internal/walletModule/provider"
)

type WalletController struct {
	walletService *walletProvider.WalletService
}

func NewWalletController(walletService *walletProvider.WalletService) *WalletController {
	return &WalletController{
		walletService,
	}
}

func (walletController *WalletController) ReadWallet(context *gin.Context, query *walletEntity.ReadWalletRequest) {
	readWalletResult, err := walletController.walletService.ReadWallet(context)
	if err != nil {
		middleware.ResponseServerError(context, err)
		return
	}
	middleware.ResponseOK(context, readWalletResult)
}

func (walletController *WalletController) Operation(context *gin.Context, body *walletEntity.WalletOperationRequest) {
	walletController.walletService.Operation(context, body)

	middleware.ResponseOKWithoutData(context)
}
