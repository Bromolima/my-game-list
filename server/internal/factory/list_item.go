package factory

import (
	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/google/uuid"
)

func NewListItem(gameListID, gameID uuid.UUID, status string, rating float32) *entities.ListItem {
	return &entities.ListItem{
		GameListID: gameListID,
		GameID:     gameID,
		Status:     status,
		Rating:     rating,
	}
}
