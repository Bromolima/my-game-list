package repository

import (
	"context"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/game_list_repository.go -package=mocks
type GameListRepository interface {
	BaseRepository[entities.GameList, uuid.UUID]
	FindGamesByListID(ctx context.Context, listID uuid.UUID) ([]*entities.Game, error)
}

type gameListRepository struct {
	BaseRepository[entities.GameList, uuid.UUID]
	db     *gorm.DB
	logger *slog.Logger
}

func NewGameListRepository(db *gorm.DB, logger *slog.Logger) GameListRepository {
	return &gameListRepository{
		BaseRepository: NewBaseRepository[entities.GameList, uuid.UUID](db, logger),
		db:             db,
		logger:         logger.With(slog.String("gameList", "repository")),
	}
}

func (r *gameListRepository) FindGamesByListID(ctx context.Context, listID uuid.UUID) ([]*entities.Game, error) {
	log := r.logger.With(slog.String("func", "GetGamesByListID"))

	var games []*entities.Game
	if err := r.db.Joins("JOIN list_items li ON li.game_id = games.id").
		Where("li.list_id = ?", listID).
		Find(&games).Error; err != nil {
		log.Error("Failed to find games by list id", "error", err.Error())
	}

	return games, nil
}
