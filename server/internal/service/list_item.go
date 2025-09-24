package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/factory"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListItemService interface {
	AddGameToList(ctx context.Context, userID, gameListID, gameID uuid.UUID, status string, rating float32) *resterr.RestErr
	UpdateGameFromList(ctx context.Context, gameID, gameListID, userID uuid.UUID, rating float32, status string) *resterr.RestErr
	DeleteGameFromList(ctx context.Context, gameId, gameListID, userID uuid.UUID) *resterr.RestErr
}

type listItemService struct {
	listItemRepo repository.ListItemRepository
	gameRepo     repository.GameRepository
	gameListRepo repository.GameListRepository
	logger       *slog.Logger
}

func NewListItemService(listItemRepo repository.ListItemRepository, gameRepo repository.GameRepository, gameListRepo repository.GameListRepository, logger *slog.Logger) ListItemService {
	return &listItemService{
		listItemRepo: listItemRepo,
		gameRepo:     gameRepo,
		gameListRepo: gameListRepo,
		logger:       logger.With(slog.String("service", "listItem")),
	}
}

func (s *listItemService) AddGameToList(ctx context.Context, userID, gameID, gameListID uuid.UUID, status string, rating float32) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "AddGameToList"))

	_, err := s.gameRepo.Find(ctx, gameID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Game not found in database")
			return resterr.NewNotFoundError("The specified game was not found")
		}

		log.Error("Failed to find game in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An unexpected error occurred while retrieving the game")
	}

	gameList, err := s.gameListRepo.Find(ctx, gameListID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Game list not found in database")
			return resterr.NewNotFoundError("The specified game list was not found")
		}

		log.Error("Failed to find game list in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An unexpected error occurred while retrieving the game list")
	}

	if gameList.UserID != userID {
		log.Warn("Unauthorized attempt to modify game list")
		return resterr.NewForbiddenError("You do not have permission to modify this list")
	}

	listItem := factory.NewListItem(gameListID, gameID, status, rating)
	if err := s.listItemRepo.Create(ctx, listItem); err != nil {
		log.Error("Failed to add game to list in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Could not add the game to the list")
	}

	return nil
}

func (s *listItemService) UpdateGameFromList(ctx context.Context, gameID, gameListID, userID uuid.UUID, rating float32, status string) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "UpdateGameFromList"))

	_, err := s.gameRepo.Find(ctx, gameID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Game not found in database")
			return resterr.NewNotFoundError("The specified game was not found")
		}

		log.Error("Failed to find game in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An unexpected error occurred while retrieving the game")
	}

	gameList, err := s.gameListRepo.Find(ctx, gameListID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Game list not found in database")
			return resterr.NewNotFoundError("The specified game list was not found")
		}

		log.Error("Failed to find game list in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An unexpected error occurred while retrieving the game list")
	}

	if gameList.UserID != userID {
		log.Warn("Unauthorized attempt to modify game list")
		return resterr.NewForbiddenError("You do not have permission to modify this list")
	}

	if err := s.listItemRepo.Update(ctx, gameID, gameListID, rating, status); err != nil {
		log.Error("Failed to update game in list in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Could not add the game to the list")
	}

	return nil
}

func (s *listItemService) DeleteGameFromList(ctx context.Context, gameID, gameListID, userID uuid.UUID) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "DeleteGameFromList"))

	_, err := s.gameRepo.Find(ctx, gameID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Game not found in database")
			return resterr.NewNotFoundError("The specified game was not found")
		}

		log.Error("Failed to find game in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An unexpected error occurred while retrieving the game")
	}

	gameList, err := s.gameListRepo.Find(ctx, gameListID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Game list not found in database")
			return resterr.NewNotFoundError("The specified game list was not found")
		}

		log.Error("Failed to find game list in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An unexpected error occurred while retrieving the game list")
	}

	if gameList.UserID != userID {
		log.Warn("Unauthorized attempt to modify game list")
		return resterr.NewForbiddenError("You do not have permission to modify this list")
	}

	if err := s.listItemRepo.Delete(ctx, gameID, gameListID); err != nil {
		log.Error("Failed to delete game from list in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("Could not add the game to the list")
	}

	return nil
}
