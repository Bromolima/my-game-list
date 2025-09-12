package repository

import (
	"context"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(context.Context, *models.User) error
	FindUser(context.Context, string) (*models.User, error)
	FindUserByEmail(context.Context, string) (*models.User, error)
	UpdateUser(context.Context, *models.User) error
	DeleteUser(context.Context, string) error
}

type userRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger.With(slog.String("repository", "user")),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	log := r.logger.With(slog.String("func", "CreateUser"))

	user.ID = uuid.NewString()
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *userRepository) FindUser(ctx context.Context, id string) (*models.User, error) {
	log := r.logger.With(slog.String("func", "FindUser"))

	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		log.Error("failed to find user", slog.String("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	log := r.logger.With(slog.String("func", "FindUserByEmail"))

	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		log.Error("failed to find user by email", slog.String("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	log := r.logger.With(slog.String("func", "UpdateUser"))

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Error("failed to update user", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	log := r.logger.With(slog.String("func", "DeleteUser"))

	uniqueID, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse uuid", slog.String("error", err.Error()))
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.User{}, uniqueID).Error; err != nil {
		log.Error("failed to delete user", slog.String("error", err.Error()))
	}

	return nil
}
