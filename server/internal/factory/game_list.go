package factory

import (
	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/http/dto"
	"github.com/google/uuid"
)

const (
	DefaultListName = "Watchlist"
)

func NewGameList(userID uuid.UUID, name string, isPublic, isDefault bool) *entities.GameList {
	return &entities.GameList{
		ID:        uuid.New(),
		Name:      name,
		IsPublic:  isDefault,
		IsDefault: isDefault,
		UserID:    userID,
	}
}

func NewGameListUpdate(gameListID, userID uuid.UUID, name string, isPublic bool) *entities.GameList {
	return &entities.GameList{
		ID:       gameListID,
		UserID:   userID,
		Name:     name,
		IsPublic: isPublic,
	}
}

func NewDefaultGameList(userID uuid.UUID) *entities.GameList {
	return &entities.GameList{
		ID:        uuid.New(),
		Name:      DefaultListName,
		IsPublic:  true,
		IsDefault: true,
		UserID:    userID,
	}
}

func NewResponseFromGameList(gameList *entities.GameList) *dto.GameListResponse {
	return &dto.GameListResponse{
		Name: gameList.Name,
	}
}
