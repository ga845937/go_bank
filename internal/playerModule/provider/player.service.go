package playerProvider

import (
	"context"
	"go_bank/internal/db/postgres"
	"go_bank/internal/db/postgres/model/player"
	playerEntity "go_bank/internal/playerModule/entity"
	"log"
)

type PlayerService struct {
	postgresModel *postgres.PostgresModel
}

func NewPlayerService(postgresModel *postgres.PostgresModel) *PlayerService {
	return &PlayerService{
		postgresModel,
	}
}

func (playerService *PlayerService) CreatePlayer(context context.Context) (string, error) {
	createPlayerResult, err := playerService.postgresModel.Player().Create().SetID("test@gmail").SetName("kuo").SetStatus(player.StatusACTIVE).Save(context)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(createPlayerResult)

	return "1", nil
}

func (playerService *PlayerService) ReadPlayer(context context.Context, query *playerEntity.ReadPlayerRequest) (*playerEntity.ReadPlayerResponse, error) {
	readPlayerResult, err := playerService.postgresModel.Player().Get(context, query.Email)
	if err != nil {
		return nil, err
	}

	data := &playerEntity.ReadPlayerResponse{
		Email:  readPlayerResult.ID,
		Name:   readPlayerResult.Name,
		Status: readPlayerResult.Status,
	}

	return data, nil
}

// func (us *PlayerService) CreatePlayer(req CreatePlayerRequest) (Player, error) {
// 	if req.Name == "" {
// 		return Player{}, errors.New("name is required")
// 	}
