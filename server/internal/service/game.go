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

type GameService interface {
	CreateGame(ctx context.Context, name, genre, developer, description, imageURL string) *resterr.RestErr
	FindGame(ctx context.Context, id uuid.UUID) (*entities.Game, *resterr.RestErr)
	SearchGames(ctx context.Context, page *entities.Page[entities.Game], query string) (*entities.Page[entities.Game], *resterr.RestErr)
	UpdateGame(ctx context.Context, id uuid.UUID, name, genre, developer, description, imageURL string) *resterr.RestErr
	DeleteGame(ctx context.Context, id string) *resterr.RestErr
}

type gameService struct {
	repository repository.GameRepository
	logger     *slog.Logger
}

func NewGameService(repository repository.GameRepository, logger *slog.Logger) GameService {
	return &gameService{
		repository: repository,
		logger:     logger.With(slog.String("service", "game")),
	}
}

func (s *gameService) CreateGame(ctx context.Context, name, genre, developer, description, imageURL string) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "CreateGame"))
	game := factory.NewGame(name, genre, developer, description, imageURL)
	if err := s.repository.Create(ctx, game); err != nil {
		log.Error("Failed to create game in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while creating the game")
	}

	return nil
}

func (s *gameService) FindGame(ctx context.Context, id uuid.UUID) (*entities.Game, *resterr.RestErr) {
	log := s.logger.With(slog.String("func", "FindGame"))

	game, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("The requested game was not found")
			return nil, resterr.NewNotFoundError("The requested game was not found")
		}

		log.Error("Failed to find game in database", slog.String("error", err.Error()))
		return nil, resterr.NewInternalServerErr("An error occurred while finding the game")
	}

	return game, nil
}

func (s *gameService) SearchGames(ctx context.Context, page *entities.Page[entities.Game], query string) (*entities.Page[entities.Game], *resterr.RestErr) {
	log := s.logger.With(slog.String("func", "SearchGames"))

	page, err := s.repository.Search(ctx, page, query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No games were found for the given query")
			return nil, resterr.NewNotFoundError("No games were found for the given query")
		}

		log.Error("Failed to search for games in database", slog.String("error", err.Error()))
		return nil, resterr.NewInternalServerErr("An error occurred while searching for games")
	}

	return page, nil
}

func (s *gameService) UpdateGame(ctx context.Context, id uuid.UUID, name, genre, developer, description, imageURL string) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "UpdateGame"))

	_, err := s.repository.Find(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("The requested game was not found")
			return resterr.NewNotFoundError("The requested game was not found")
		}

		log.Error("Failed to find game in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while finding the game")
	}
	game := factory.NewGameUpdate(id, name, genre, developer, description, imageURL)
	if err := s.repository.Update(ctx, game); err != nil {
		log.Error("Failed to update game in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while updating the game")
	}

	return nil
}

func (s *gameService) DeleteGame(ctx context.Context, id string) *resterr.RestErr {
	log := s.logger.With(slog.String("func", "DeleteGame"))

	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Error("Failed to parse game ID", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while parsing the UUID")
	}

	_, err = s.repository.Find(ctx, uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("The requested game was not found")
			return resterr.NewNotFoundError("The requested game was not found")
		}

		log.Error("Failed to find game in database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while finding the game")
	}

	if err := s.repository.Delete(ctx, uuid); err != nil {
		log.Error("Failed to delete game from database", slog.String("error", err.Error()))
		return resterr.NewInternalServerErr("An error occurred while deleting the game")
	}

	return nil
}
