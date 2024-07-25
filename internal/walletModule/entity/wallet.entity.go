package walletType

import (
	_ "github.com/go-playground/validator/v10"
)

type ReadWalletRequest struct {
	ID string `form:"id" binding:"required"`
}

type ReadWalletResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Balance int    `json:"balance"`
}

type WalletOperationType string

const (
	Deposit    WalletOperationType = "Deposit"
	Withdrawal WalletOperationType = "Withdrawal"
)

type WalletOperation struct {
	ID             string
	Operation      WalletOperationType
	Amount         int
	RequestTraceID string
}

type WalletOperationRequest struct {
	ID        string              `form:"id" binding:"required"`
	Operation WalletOperationType `form:"operation" binding:"oneof=Deposit Withdrawal"`
	Amount    int                 `form:"amount" binding:"gte=1"`
}
