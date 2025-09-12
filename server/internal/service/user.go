package service

import (
	"context"
	"errors"
	"log/slog"

	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/models"
	"github.com/Bromolima/my-game-list/internal/repository"
	"github.com/Bromolima/my-game-list/internal/token"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(context.Context, *models.User) *resterr.RestErr
	FindUser(context.Context, string) (*models.User, *resterr.RestErr)
	UpdateUser(context.Context, *models.User) *resterr.RestErr
	DeleteUser(context.Context, string) *resterr.RestErr
	Login(context.Context, *models.User) (*models.User, string, *resterr.RestErr)
}

type userService struct {
	repository   repository.UserRepository
	tokenService token.JwtService
	logger       *slog.Logger
}

func NewUserService(repository repository.UserRepository, tokenService token.JwtService, logger *slog.Logger) UserService {
	return &userService{
		repository:   repository,
		tokenService: tokenService,
		logger:       logger.With(slog.String("service", "user")),
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "CreateUser"))

	userExists, err := s.repository.FindUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("failed to find user by email", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Error trying to find user by Email")
	}

	if userExists != nil {
		log.Warn("email already registered")
		return resterr.NewBadRequestError("Email is already registeresd in another account")
	}

	user.Password, err = models.HashPassword(user.Password)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Failed to hash password")
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Error trying to create user")
	}

	return nil
}

func (s *userService) FindUser(ctx context.Context, id string) (*models.User, *resterr.RestErr) {
	log := s.logger.With(slog.String("func", "FindUser"))

	user, err := s.repository.FindUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("user not found")
			return nil, resterr.NewNotFoundError("User not found")
		}

		log.Error("failed to find user", slog.String("error", err.Error()))
		return nil, resterr.NewInternalServerErr("Error trying to find user")
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "UpdateUser"))

	_, err := s.repository.FindUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("user not found")
			return resterr.NewNotFoundError("User not found")
		}

		log.Error("failed to find user by email", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Error trying to find user")
	}

	if err := s.repository.UpdateUser(ctx, user); err != nil {
		log.Error("failed to update user", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Error trying to update user")
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, ID string) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "DeleteUser"))

	_, err := s.repository.FindUser(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("user not found")
			return resterr.NewNotFoundError("User not found")
		}

		log.Error("failed to find user", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Error trying to find user")
	}

	if err := s.repository.DeleteUser(ctx, ID); err != nil {
		log.Error("failed to delete user", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Error trying to update user")
	}

	return nil
}
