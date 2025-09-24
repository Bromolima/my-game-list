package repository

import (
	"context"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockgen -source=user.go -destination=../../mocks/user_repository.go -package=mocks
type UserRepository interface {
	BaseRepository[entities.User, uuid.UUID]
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Search(ctx context.Context, page *entities.Page[entities.User], query string) (*entities.Page[entities.User], error)
}

type userRepository struct {
	BaseRepository[entities.User, uuid.UUID]
	db             *gorm.DB
	pageRepository PageRepository[entities.User]
	logger         *slog.Logger
}

func NewUserRepository(db *gorm.DB, pageRepository PageRepository[entities.User], logger *slog.Logger) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[entities.User, uuid.UUID](db, logger),
		db:             db,
		pageRepository: pageRepository,
		logger:         logger.With(slog.String("repository", "user")),
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	log := r.logger.With(slog.String("func", "FindByEmail"))

	var user entities.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		log.Error("Failed to find user by email in database", slog.String("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Search(ctx context.Context, page *entities.Page[entities.User], query string) (*entities.Page[entities.User], error) {
	log := r.logger.With(slog.String("func", "Search"))

	search := "%" + query + "%"
	var data []entities.User
	if err := r.db.WithContext(ctx).
		Scopes(r.pageRepository.Paginate(&entities.User{}, page)).
		Where("username LIKE ?", search).
		Find(&data).Error; err != nil {
		log.Error("Failed to search users in database", slog.String("error", err.Error()))
		return nil, err
	}

	page.Data = data
	return page, nil
}
