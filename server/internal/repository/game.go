package repository

import (
	"context"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockgen -source=game.go -destination=../../mocks/game_repository.go -package=mocks
type GameRepository interface {
	BaseRepository[entities.Game, uuid.UUID]
	Search(ctx context.Context, page *entities.Page[entities.Game], query string) (*entities.Page[entities.Game], error)
}

type gameRepository struct {
	BaseRepository[entities.Game, uuid.UUID]
	db             *gorm.DB
	pageRepository PageRepository[entities.Game]
	logger         *slog.Logger
}

func NewGameRepository(db *gorm.DB, logger *slog.Logger, pageRepository PageRepository[entities.Game]) GameRepository {
	return &gameRepository{
		BaseRepository: NewBaseRepository[entities.Game, uuid.UUID](db, logger),
		db:             db,
		pageRepository: pageRepository,
		logger:         logger.With(slog.String("game", "repository")),
	}
}

func (r *gameRepository) Search(ctx context.Context, page *entities.Page[entities.Game], query string) (*entities.Page[entities.Game], error) {
	log := r.logger.With(slog.String("func", "Search"))

	search := "%" + query + "%"
	var data []entities.Game
	if err := r.db.WithContext(ctx).
		Scopes(r.pageRepository.Paginate(&entities.Game{}, page)).
		Where("name ILIKE ?", search).
		Find(&data).Error; err != nil {
		log.Error("Failed to search games in database", slog.String("error", err.Error()))
		return nil, err
	}

	page.Data = data
	return page, nil
}
