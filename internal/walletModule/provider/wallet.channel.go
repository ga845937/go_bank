package walletProvider

import (
	"context"
	"go_bank/internal/db/postgres"
	walletEntity "go_bank/internal/walletModule/entity"

	"go.elastic.co/apm/v2"
)

type WalletChannel struct {
	channel       chan walletEntity.WalletOperation
	postgresModel *postgres.PostgresModel
}

func NewWalletChannel(postgresModel *postgres.PostgresModel) *WalletChannel {
	return &WalletChannel{
		channel:       make(chan walletEntity.WalletOperation, 100),
		postgresModel: postgresModel,
	}
}

func (walletChannel *WalletChannel) SendOperation(walletOperation walletEntity.WalletOperation) {
	walletChannel.channel <- walletOperation
}

func (walletChannel *WalletChannel) Start() {
	go func() {
		for walletOperation := range walletChannel.channel {
			walletChannel.handleOperation(walletOperation)
		}
	}()
}

func (walletChannel *WalletChannel) handleOperation(walletOperation walletEntity.WalletOperation) {
	transaction := apm.DefaultTracer().StartTransaction("walletChannel", "request")
	defer transaction.End()

	transaction.Context.SetLabel("requestTraceID", walletOperation.RequestTraceID)

	ctx := apm.ContextWithTransaction(context.Background(), transaction)

	tx, err := walletChannel.postgresModel.Tx(ctx)
	if err != nil {
		transaction.Context.SetLabel("error", err.Error())
		tx.Rollback()
		return
	}

	if walletOperation.Operation == walletEntity.Withdrawal {
		wallet, err := tx.Wallet.Get(ctx, walletOperation.ID)
		if err != nil {
			transaction.Context.SetLabel("error", err.Error())
			tx.Rollback()
			return
		}

		if wallet.Balance < walletOperation.Amount {
			transaction.Context.SetLabel("error", "insufficient balance")
			tx.Rollback()
			return
		}
	}

	amount := walletOperation.Amount
	if walletOperation.Operation == walletEntity.Withdrawal {
		amount = -amount
	}

	walletResult, err := tx.Wallet.UpdateOneID(walletOperation.ID).AddBalance(amount).Save(ctx)
	if err != nil {
		transaction.Context.SetLabel("error", err.Error())
		tx.Rollback()
		return
	}

	operationTraceID := transaction.TraceContext().Trace.String()
	_, err = tx.Record.Create().SetID(walletOperation.RequestTraceID).SetUserEmail(walletResult.UserEmail).SetWalletID(walletOperation.ID).SetOperationTraceID(operationTraceID).SetAmount(amount).SetBalance(walletResult.Balance).Save(ctx)
	if err != nil {
		transaction.Context.SetLabel("error", err.Error())
		tx.Rollback()
		return
	}

	tx.Commit()
}
