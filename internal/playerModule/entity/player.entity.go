package playerType

import (
	"go_bank/internal/db/postgres/model/player"

	_ "github.com/go-playground/validator/v10"
)

type CreatePlayerRequest struct {
	Email  string `json:"email" binding:"email"`
	Name   string `json:"name" binding:"required"`
	Status string `json:"status" binding:"oneof=INACTIVE ACTIVE BANNED"`
}

type ReadPlayerRequest struct {
	Email string `form:"email" binding:"email"`
}

type ReadPlayerResponse struct {
	Email  string        `json:"email"`
	Name   string        `json:"name"`
	Status player.Status `json:"status"`
}
