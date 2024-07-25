package route

import (
	"os"

	"go_bank/internal/db/postgres"
	playerModule "go_bank/internal/playerModule"
	playerProvider "go_bank/internal/playerModule/provider"
	walletModule "go_bank/internal/walletModule"
	walletProvider "go_bank/internal/walletModule/provider"

	"go_bank/internal/middleware"

	"go.elastic.co/apm/module/apmgin/v2"

	"github.com/gin-gonic/gin"
)

func InitializeRoute() {
	postgresModel := postgres.InitializePostgres()

	engine := gin.New()
	engine.Use(apmgin.Middleware(engine), middleware.Logger())

	walletChannel := walletProvider.NewWalletChannel(postgresModel)
	walletChannel.Start()
	walletService := walletProvider.NewWalletService(postgresModel, walletChannel)
	walletController := walletModule.NewWalletController(walletService)
	walletRoute := walletModule.NewWalletRoute(walletController)
	walletRoute.RegisterRoute(engine.Group("/wallet"))

	playerService := playerProvider.NewPlayerService(postgresModel)
	playerController := playerModule.NewPlayerController(playerService)
	playerRoute := playerModule.NewPlayerRoute(playerController)
	playerRoute.RegisterRoute(engine.Group("/player"))

	if err := engine.Run(os.Getenv("HTTP_PORT")); err != nil {
		panic(err)
	}
}
