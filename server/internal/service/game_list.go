package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/factory"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameListService interface {
	CreateGameList(ctx context.Context, userID uuid.UUID, name string, isPublic, isDefault bool) *resterr.RestErr
	FindGamesFromList(ctx context.Context, gameListID uuid.UUID) ([]*entities.Game, *resterr.RestErr)
	UpdateGameList(ctx context.Context, userID, gameListID uuid.UUID, name string, isPublic bool) *resterr.RestErr
	DeleteGameList(ctx context.Context, gameListID, userID uuid.UUID) *resterr.RestErr
}

type gameListService struct {
	gameListRepo repository.GameListRepository
	userRepo     repository.UserRepository
	logger       *slog.Logger
}

func NewGameListService(gameListRepo repository.GameListRepository, userRepo repository.UserRepository, logger *slog.Logger) GameListService {
	return &gameListService{
		gameListRepo: gameListRepo,
		userRepo:     userRepo,
		logger:       logger.With(slog.String("service", "gameList")),
	}
}

func (s *gameListService) CreateGameList(ctx context.Context, userID uuid.UUID, name string, isPublic, isDefault bool) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "CreateGameList"))

	_, err := s.userRepo.Find(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found in database")
			restErr := resterr.NewNotFoundError("User does not exists")
			return restErr
		}

		log.Error("Failed to find user in database", slog.String("error", err.Error()))
		restErr := resterr.NewInternalServerErr("Failed to search for user due to internal error")
		return restErr
	}

	gameList := factory.NewGameList(userID, name, isPublic, isDefault)
	if err := s.gameListRepo.Create(ctx, gameList); err != nil {
		log.Error("Failed to create game list in database", slog.String("error", err.Error()))
		restErr := resterr.NewInternalServerErr("Failed to create list due to internal error")
		return restErr
	}

	return nil
}

func (s *gameListService) FindGamesFromList(ctx context.Context, gameListID uuid.UUID) ([]*entities.Game, *resterr.RestErr) {
	log := s.logger.With(slog.String("func", "FindGamesByList"))

	games, err := s.gameListRepo.FindGamesByListID(ctx, gameListID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Game list not found, returning empty list", "gameListID", gameListID.String())
			return []*entities.Game{}, nil
		}
		log.Error("Failed to find games in list", "error", slog.String("error", err.Error()))
		restErr := resterr.NewInternalServerErr("Failed to find games in list")
		return nil, restErr
	}

	return games, nil
}

func (s *gameListService) UpdateGameList(ctx context.Context, userID, gameListID uuid.UUID, name string, isPublic bool) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "UpdateGameList"))

	_, err := s.userRepo.Find(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found in database")
			restErr := resterr.NewNotFoundError("User does not exist")
			return restErr
		}

		log.Error("Failed to find user in database", slog.String("error", err.Error()))
		restErr := resterr.NewInternalServerErr("Failed to search for user due to internal error")
		return restErr
	}

	gameList := factory.NewGameListUpdate(gameListID, userID, name, isPublic)
	if err := s.gameListRepo.Update(ctx, gameList); err != nil {
		log.Error("Failed to update game list in database", slog.String("error", err.Error()))
		restErr := resterr.NewInternalServerErr("Failed to create list due to internal error")
		return restErr
	}

	return nil
}

func (s *gameListService) DeleteGameList(ctx context.Context, gameListID, userID uuid.UUID) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "DeleteGameList"))

	_, err := s.userRepo.Find(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found in database")
			restErr := resterr.NewNotFoundError("User does not exists")
			return restErr
		}

		log.Error("Failed to find user in database", slog.String("error", err.Error()))
		restErr := resterr.NewInternalServerErr("Failed to search for user due to internal error")
		return restErr
	}

	_, err = s.gameListRepo.Find(ctx, gameListID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Game list not found in database")
			restErr := resterr.NewNotFoundError("Game list does not exists")
			return restErr
		}

		log.Error("Failed to find game list in database", slog.String("error", err.Error()))
		restErr := resterr.NewInternalServerErr("Failed to search for game list due to internal error")
		return restErr
	}

	if err := s.gameListRepo.Delete(ctx, gameListID); err != nil {
		log.Error("Failed to delete game list from database", slog.String("error", err.Error()))
		restErr := resterr.NewInternalServerErr("Failed to delete list due to internal error")
		return restErr
	}

	return nil
}
