package walletProvider

import (
	"go_bank/internal/db/postgres"
	walletEntity "go_bank/internal/walletModule/entity"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/v2"
)

type WalletService struct {
	postgresModel *postgres.PostgresModel
	walletChannel *WalletChannel
}

func NewWalletService(postgresModel *postgres.PostgresModel, walletChannel *WalletChannel) *WalletService {
	return &WalletService{
		postgresModel: postgresModel,
		walletChannel: walletChannel,
	}
}

func (walletService *WalletService) ReadWallet(context *gin.Context) (*walletEntity.ReadWalletResponse, error) {
	var traceID string
	tx := apm.TransactionFromContext(context.Request.Context())
	if tx != nil {
		traceID = tx.TraceContext().Trace.String()
	}

	readWalletResult, err := walletService.postgresModel.Wallet().Get(context, traceID)
	if err != nil {
		return nil, err
	}

	data := &walletEntity.ReadWalletResponse{
		ID:      traceID,
		Email:   readWalletResult.UserEmail,
		Balance: readWalletResult.Balance,
	}

	return data, nil
}

func (walletService *WalletService) Operation(context *gin.Context, walletOperation *walletEntity.WalletOperationRequest) {
	var traceID string
	tx := apm.TransactionFromContext(context.Request.Context())
	if tx != nil {
		traceID = tx.TraceContext().Trace.String()
	}

	operationData := walletEntity.WalletOperation{
		ID:             walletOperation.ID,
		Operation:      walletOperation.Operation,
		Amount:         walletOperation.Amount,
		RequestTraceID: traceID,
	}

	walletService.walletChannel.SendOperation(operationData)
}
