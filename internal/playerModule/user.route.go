package playerModule

import (
	"github.com/gin-gonic/gin"

	"go_bank/internal/middleware"
)

type PlayerRoute struct {
	playerController *PlayerController
}

func NewPlayerRoute(playerController *PlayerController) *PlayerRoute {
	return &PlayerRoute{
		playerController,
	}
}

func (playerRoute *PlayerRoute) RegisterRoute(playerRouteGroup *gin.RouterGroup) {
	playerRouteGroup.POST("/", middleware.Validator(playerRoute.playerController.CreatePlayer))
	playerRouteGroup.GET("/", middleware.Validator(playerRoute.playerController.ReadPlayer))
}
