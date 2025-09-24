package repository

import (
	"math"

	"github.com/Bromolima/my-game-list/internal/entities"
	"gorm.io/gorm"
)

type PageRepository[T any] interface {
	Paginate(data *T, page *entities.Page[T]) func(db *gorm.DB) *gorm.DB
}

type pageRepository[T any] struct {
	db *gorm.DB
}

func NewPageRepository[T any](db *gorm.DB) PageRepository[T] {
	return &pageRepository[T]{
		db: db,
	}
}

func (r *pageRepository[T]) Paginate(data *T, page *entities.Page[T]) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	r.db.Model(data).Count(&totalRows)

	page.TotalItems = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(page.Limit)))
	page.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(page.Offset).Limit(page.Limit)
	}
}
