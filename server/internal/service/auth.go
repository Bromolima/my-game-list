package service

import (
	"errors"
	"log/slog"

	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/models"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

func (s *userService) Login(ctx context.Context, user *models.User) (*models.User, string, *resterr.RestErr) {
	log := s.logger.With(slog.String("func", "Login"))

	userExists, err := s.repository.FindUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("failed to find user by email", slog.String("error", err.Error()))
		return nil, "", resterr.NewInternalServerErr("Failed to find user")
	}

	if userExists == nil {
		log.Warn("user not found")
		return nil, "", resterr.NewNotFoundError("User not found")
	}

	if !models.CheckPassword(userExists.Password, user.Password) {
		log.Warn("invalid credentials")
		return nil, "", resterr.NewUnauthorizedError("User with wrong credential")
	}

	token, err := s.tokenService.GenerateToken(userExists)
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
		return nil, "", resterr.NewInternalServerErr("Failed to generate token")
	}

	return userExists, token, nil
}
