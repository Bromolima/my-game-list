package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/factory"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/repository"
	"github.com/Bromolima/my-game-list/internal/security"
	"github.com/Bromolima/my-game-list/internal/token"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(ctx context.Context, email, password, username, avatarURL string) *resterr.RestErr
	FindUser(ctx context.Context, id string) (*entities.User, *resterr.RestErr)
	SearchUsers(ctx context.Context, page *entities.Page[entities.User], query string) (*entities.Page[entities.User], *resterr.RestErr)
	UpdateUser(ctx context.Context, id uuid.UUID, email, password, username, avatarURL string) *resterr.RestErr
	DeleteUser(ctx context.Context, id string) *resterr.RestErr
	Login(ctx context.Context, email, password string) (string, *resterr.RestErr)
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

func (s *userService) RegisterUser(ctx context.Context, email, password, username, avatarURL string) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "RegisterUser"))

	userExists, err := s.repository.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("Failed to find user by email", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while finding the user by email")
	}

	if userExists != nil {
		log.Warn("The provided email is already registered")
		return resterr.NewBadRequestError("The provided email is already registered")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		log.Error("Failed to hash password", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while hashing the password")
	}

	user := factory.NewUser(email, hashedPassword, username, avatarURL, entities.RoleUserID)
	if err := s.repository.Create(ctx, user); err != nil {
		log.Error("Failed to create user in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while creating the user")
	}

	return nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, *resterr.RestErr) {
	log := s.logger.With(slog.String("func", "Login"))

	userExists, err := s.repository.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("Failed to find user by email", slog.String("error", err.Error()))
		return "", resterr.NewInternalServerErr("An error occurred while finding the user")
	}

	if userExists == nil {
		log.Warn("The requested user was not found")
		return "", resterr.NewNotFoundError("The requested user was not found")
	}

	if !security.CheckPassword(userExists.Password, password) {
		log.Warn("Invalid credentials provided")
		return "", resterr.NewUnauthorizedError("Invalid credentials provided")
	}

	token, err := s.tokenService.GenerateToken(userExists)
	if err != nil {
		log.Error("Failed to generate token", slog.String("error", err.Error()))
		return "", resterr.NewInternalServerErr("An error occurred while generating the token")
	}

	return token, nil
}

func (s *userService) FindUser(ctx context.Context, id string) (*entities.User, *resterr.RestErr) {
	log := s.logger.With(slog.String("func", "FindUser"))

	user, err := s.repository.Find(ctx, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("The requested user was not found")
			return nil, resterr.NewNotFoundError("The requested user was not found")
		}

		log.Error("Failed to find user in database", slog.String("error", err.Error()))
		return nil, resterr.NewInternalServerErr("An error occurred while finding the user")
	}

	return user, nil
}

func (s *userService) SearchUsers(ctx context.Context, page *entities.Page[entities.User], query string) (*entities.Page[entities.User], *resterr.RestErr) {
	log := s.logger.With(slog.String("func", "SearchUsers"))

	page, err := s.repository.Search(ctx, page, query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No users were found for the given query")
			return nil, resterr.NewNotFoundError("No users were found for the given query")
		}

		log.Error("Failed to search for users in database", slog.String("error", err.Error()))
		return nil, resterr.NewInternalServerErr("An error occurred while searching for users")
	}

	return page, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, email, password, username, avatarURL string) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "UpdateUser"))

	_, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("The requested user was not found")
			return resterr.NewNotFoundError("The requested user was not found")
		}

		log.Error("Failed to find user by email", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while finding the user")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		log.Error("Failed to hash password", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while hashing the password")
	}

	user := factory.NewUserUpdate(id, email, hashedPassword, username, avatarURL)
	if err := s.repository.Update(ctx, user); err != nil {
		log.Error("Failed to update user in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while updating the user")
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "DeleteUser"))

	uniqueID := uuid.MustParse(id)
	_, err := s.repository.Find(ctx, uniqueID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("The requested user was not found")
			return resterr.NewNotFoundError("The requested user was not found")
		}

		log.Error("Failed to find user in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while finding the user")
	}

	if err := s.repository.Delete(ctx, uniqueID); err != nil {
		log.Error("Failed to delete user from database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while deleting the user")
	}

	return nil
}
