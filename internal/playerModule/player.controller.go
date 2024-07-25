package playerModule

import (
	"github.com/gin-gonic/gin"

	middleware "go_bank/internal/middleware"
	playerEntity "go_bank/internal/playerModule/entity"
	playerProvider "go_bank/internal/playerModule/provider"
)

type PlayerController struct {
	playerService *playerProvider.PlayerService
}

func NewPlayerController(playerService *playerProvider.PlayerService) *PlayerController {
	return &PlayerController{
		playerService,
	}
}

func (playerController *PlayerController) ReadPlayer(context *gin.Context, query *playerEntity.ReadPlayerRequest) {
	readPlayerResult, err := playerController.playerService.ReadPlayer(context.Request.Context(), query)
	if err != nil {
		middleware.ResponseServerError(context, err)
		return
	}
	middleware.ResponseOK(context, readPlayerResult)
}

func (playerController *PlayerController) CreatePlayer(context *gin.Context, query *playerEntity.ReadPlayerRequest) {
	readPlayerResult, err := playerController.playerService.ReadPlayer(context.Request.Context(), query)
	if err != nil {
		middleware.ResponseServerError(context, err)
		return
	}
	middleware.ResponseOK(context, readPlayerResult)
}

// func (uc *PlayerController) CreatePlayer(c *gin.Context) {
// 	var newPlayer playerEntity.CreatePlayerRequest
// 	if err := context.ShouldBindJSON(&newPlayer); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	player, err := ucontext.playerService.CreatePlayer(newPlayer)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	context.JSON(http.StatusCreated, player)
// }
