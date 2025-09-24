package repository

import (
	"context"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockgen -source=list_item.go -destination=../../mocks/list_item_repository.go -package=mocks
type ListItemRepository interface {
	Create(ctx context.Context, listItem *entities.ListItem) error
	Find(ctx context.Context, gameID, gameListID uuid.UUID) (*entities.ListItem, error)
	Update(ctx context.Context, gameID uuid.UUID, gameListID uuid.UUID, rating float32, status string) error
	Delete(ctx context.Context, gameID, gameListID uuid.UUID) error
}

type listItemRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewListItemRepository(db *gorm.DB, logger *slog.Logger) ListItemRepository {
	return &listItemRepository{
		db:     db,
		logger: logger.With(slog.String("list_item", "repository")),
	}
}

func (r *listItemRepository) Create(ctx context.Context, listItem *entities.ListItem) error {
	log := r.logger.With(slog.String("func", "Create"))
	if err := r.db.WithContext(ctx).Create(listItem).Error; err != nil {
		log.Error("Failed to create list item in database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *listItemRepository) Find(ctx context.Context, gameID, gameListID uuid.UUID) (*entities.ListItem, error) {
	log := r.logger.With(slog.String("func", "Find"))
	var listItem entities.ListItem
	if err := r.db.WithContext(ctx).Where("game_id = ? AND game_list_id = ?", gameID, gameListID).First(&listItem).Error; err != nil {
		log.Error("Failed to find list item in database", slog.String("error", err.Error()))
		return nil, err
	}
	return &listItem, nil
}

func (r *listItemRepository) Update(ctx context.Context, gameID uuid.UUID, gameListID uuid.UUID, rating float32, status string) error {
	log := r.logger.With(slog.String("func", "Update"))
	if err := r.db.WithContext(ctx).Where("game_id = ? and game_list_id = ?", gameID, gameListID).Updates(&entities.ListItem{
		Rating: rating,
		Status: status,
	}).Error; err != nil {
		log.Error("Failed to update list item in database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *listItemRepository) Delete(ctx context.Context, gameID, gameListID uuid.UUID) error {
	log := r.logger.With(slog.String("func", "Delete"))
	if err := r.db.WithContext(ctx).
		Where("game_id = ? AND game_list_id = ?", gameID, gameListID).
		Delete(&entities.ListItem{}).Error; err != nil {
		log.Error("Failed to delete list item from database", slog.String("error", err.Error()))
		return err
	}
	return nil
}
