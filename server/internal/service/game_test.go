package service_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/service"
	"github.com/Bromolima/my-game-list/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGameService_CreateGame(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameRepository := mocks.NewMockGameRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameService := service.NewGameService(gameRepository, logger)

	ctx := context.Background()
	name := "The Witcher 3"
	genre := "RPG"
	developer := "CD Projekt Red"
	description := "A story-driven, next-generation open world role-playing game."
	imageURL := "http://example.com/witcher3.png"

	t.Run("should create game successfully", func(t *testing.T) {
		gameRepository.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		err := gameService.CreateGame(ctx, name, genre, developer, description, imageURL)

		assert.Nil(t, err)
	})

	t.Run("should return error when Create fails", func(t *testing.T) {
		gameRepository.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("database error"))

		err := gameService.CreateGame(ctx, name, genre, developer, description, imageURL)

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while creating the game", err.Message)
	})
}

func TestGameService_FindGame(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameRepository := mocks.NewMockGameRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameService := service.NewGameService(gameRepository, logger)

	ctx := context.Background()
	gameID := uuid.New()
	game := &entities.Game{ID: gameID}

	t.Run("should find game successfully", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(game, nil)

		foundGame, err := gameService.FindGame(ctx, gameID)

		assert.Nil(t, err)
		assert.Equal(t, game, foundGame)
	})

	t.Run("should return error when game is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, gorm.ErrRecordNotFound)

		foundGame, err := gameService.FindGame(ctx, gameID)

		assert.NotNil(t, err)
		assert.Nil(t, foundGame)
		assert.Equal(t, "The requested game was not found", err.Message)
	})

	t.Run("should return error when Find fails", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, errors.New("database error"))

		foundGame, err := gameService.FindGame(ctx, gameID)

		assert.NotNil(t, err)
		assert.Nil(t, foundGame)
		assert.Equal(t, "An error occurred while finding the game", err.Message)
	})
}

func TestGameService_SearchGames(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameRepository := mocks.NewMockGameRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameService := service.NewGameService(gameRepository, logger)

	ctx := context.Background()
	page := &entities.Page[entities.Game]{}
	query := "test"

	t.Run("should search games successfully", func(t *testing.T) {
		gameRepository.EXPECT().Search(ctx, page, query).Return(page, nil)

		resultPage, err := gameService.SearchGames(ctx, page, query)

		assert.Nil(t, err)
		assert.Equal(t, page, resultPage)
	})

	t.Run("should return error when no games are found", func(t *testing.T) {
		gameRepository.EXPECT().Search(ctx, page, query).Return(nil, gorm.ErrRecordNotFound)

		resultPage, err := gameService.SearchGames(ctx, page, query)

		assert.NotNil(t, err)
		assert.Nil(t, resultPage)
		assert.Equal(t, "No games were found for the given query", err.Message)
	})

	t.Run("should return error when Search fails", func(t *testing.T) {
		gameRepository.EXPECT().Search(ctx, page, query).Return(nil, errors.New("database error"))

		resultPage, err := gameService.SearchGames(ctx, page, query)

		assert.NotNil(t, err)
		assert.Nil(t, resultPage)
		assert.Equal(t, "An error occurred while searching for games", err.Message)
	})
}

func TestGameService_UpdateGame(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameRepository := mocks.NewMockGameRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameService := service.NewGameService(gameRepository, logger)

	ctx := context.Background()
	gameID := uuid.New()
	name := "The Witcher 3"
	genre := "RPG"
	developer := "CD Projekt Red"
	description := "A story-driven, next-generation open world role-playing game."
	imageURL := "http://example.com/witcher3.png"

	t.Run("should update game successfully", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameRepository.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		err := gameService.UpdateGame(ctx, gameID, name, genre, developer, description, imageURL)

		assert.Nil(t, err)
	})

	t.Run("should return error when game is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, gorm.ErrRecordNotFound)

		err := gameService.UpdateGame(ctx, gameID, name, genre, developer, description, imageURL)

		assert.NotNil(t, err)
		assert.Equal(t, "The requested game was not found", err.Message)
	})

	t.Run("should return error when Find fails", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, errors.New("database error"))

		err := gameService.UpdateGame(ctx, gameID, name, genre, developer, description, imageURL)

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while finding the game", err.Message)
	})

	t.Run("should return error when Update fails", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameRepository.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("database error"))

		err := gameService.UpdateGame(ctx, gameID, name, genre, developer, description, imageURL)

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while updating the game", err.Message)
	})
}

func TestGameService_DeleteGame(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameRepository := mocks.NewMockGameRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameService := service.NewGameService(gameRepository, logger)

	ctx := context.Background()
	gameID := uuid.New()

	t.Run("should delete game successfully", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameRepository.EXPECT().Delete(ctx, gameID).Return(nil)

		err := gameService.DeleteGame(ctx, gameID.String())

		assert.Nil(t, err)
	})

	t.Run("should return error when game is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, gorm.ErrRecordNotFound)

		err := gameService.DeleteGame(ctx, gameID.String())

		assert.NotNil(t, err)
		assert.Equal(t, "The requested game was not found", err.Message)
	})

	t.Run("should return error when Find fails", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, errors.New("database error"))

		err := gameService.DeleteGame(ctx, gameID.String())

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while finding the game", err.Message)
	})

	t.Run("should return error when Delete fails", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameRepository.EXPECT().Delete(ctx, gameID).Return(errors.New("database error"))

		err := gameService.DeleteGame(ctx, gameID.String())

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while deleting the game", err.Message)
	})
}
