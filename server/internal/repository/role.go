package repository

import (
	"context"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	HasAccess(ctx context.Context, userID uuid.UUID, accessName entities.AccessType) (error, bool)
}

type roleRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRoleRepository(db *gorm.DB, logger *slog.Logger) RoleRepository {
	return &roleRepository{
		db:     db,
		logger: logger.With(slog.String("role", "repository")),
	}
}

func (r *roleRepository) HasAccess(ctx context.Context, userID uuid.UUID, accessName entities.AccessType) (error, bool) {
	log := r.logger.With(slog.String("func", "HasAccess"))
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.Access{}).
		Joins("JOIN role_accesses ra ON ra.access_id = accesses.id").
		Joins("JOIN roles r ON r.id = ra.role_id").
		Joins("JOIN users u ON u.role_id = r.id").
		Where("u.id = ? AND accesses.access_type = ?", userID, accessName).Count(&count).Error; err != nil {
		log.Error("Failed to check user access in database", slog.String("error", err.Error()))
		return err, false
	}

	return nil, count > 0
}
