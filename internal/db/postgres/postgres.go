package postgres

import (
	"context"
	"log"
	"os"

	ent "go_bank/internal/db/postgres/model"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"go.elastic.co/apm/module/apmsql/v2"
	_ "go.elastic.co/apm/module/apmsql/v2/pq"
)

type PostgresModel struct {
	client *ent.Client
}

func InitializePostgres() *PostgresModel {
	postgresUri := os.Getenv("POSTGRES_URI")

	db, err := apmsql.Open(dialect.Postgres, postgresUri)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	dbDriver := entsql.OpenDB(dialect.Postgres, db)

	client := ent.NewClient(ent.Driver(dbDriver))

	if os.Getenv("GOLANG_ENV") == "dev" {
		client = client.Debug()
	}

	return &PostgresModel{
		client: client,
	}
}

func (model *PostgresModel) Tx(ctx context.Context) (*ent.Tx, error) {
	return model.client.Tx(ctx)
}

func (model *PostgresModel) Player() *ent.PlayerClient {
	return model.client.Player
}

func (model *PostgresModel) Wallet() *ent.WalletClient {
	return model.client.Wallet
}

func (model *PostgresModel) Record() *ent.RecordClient {
	return model.client.Record
}
