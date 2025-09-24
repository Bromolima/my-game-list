package factory

import (
	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/http/dto"
	"github.com/google/uuid"
)

func NewGame(name, genre, developer, description, imageURL string) *entities.Game {
	return &entities.Game{
		ID:          uuid.New(),
		Name:        name,
		Genre:       genre,
		Developer:   developer,
		Description: description,
		ImageURL:    imageURL,
	}
}

func NewGameUpdate(id uuid.UUID, name, genre, developer, description, imageURL string) *entities.Game {
	return &entities.Game{
		ID:          id,
		Name:        name,
		Genre:       genre,
		Developer:   developer,
		Description: description,
		ImageURL:    imageURL,
	}
}

func NewResponseFromGame(game *entities.Game) *dto.GameResponse {
	return &dto.GameResponse{
		Name:        game.Name,
		Genre:       game.Genre,
		Developer:   game.Developer,
		Rating:      game.Rating,
		Description: game.Description,
		ImageURL:    game.ImageURL,
	}
}
