package repository

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

//go:generate mockgen -source=base_repository.go -destination=../../mocks/base_repository.go -package=mocks
type BaseRepository[T any, K any] interface {
	Create(ctx context.Context, entity *T) error
	Find(ctx context.Context, id K) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id K) error
}

type baseRepository[T any, K any] struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewBaseRepository[T any, K any](db *gorm.DB, logger *slog.Logger) BaseRepository[T, K] {
	return &baseRepository[T, K]{
		db:     db,
		logger: logger.With(slog.String("base", "repository")),
	}
}

func (r *baseRepository[T, K]) Create(ctx context.Context, entity *T) error {
	log := r.logger.With(slog.String("func", "Create"))

	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		log.Error("Failed to create entity in database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *baseRepository[T, K]) Find(ctx context.Context, id K) (*T, error) {
	log := r.logger.With(slog.String("func", "Find"))
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		log.Error("Failed to find entity in database", slog.String("error", err.Error()))
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepository[T, K]) Update(ctx context.Context, entity *T) error {
	log := r.logger.With(slog.String("func", "Update"))
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		log.Error("Failed to update entity in database", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *baseRepository[T, K]) Delete(ctx context.Context, id K) error {
	log := r.logger.With(slog.String("func", "Delete"))
	var entity T
	if err := r.db.WithContext(ctx).Delete(&entity, id).Error; err != nil {
		log.Error("Failed to delete entity from database", slog.String("error", err.Error()))
		return err
	}
	return nil
}
