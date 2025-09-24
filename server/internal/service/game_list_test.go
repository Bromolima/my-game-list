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

func TestGameListService_CreateGameList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameListRepository := mocks.NewMockGameListRepository(mockCtrl)
	userRepository := mocks.NewMockUserRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameListService := service.NewGameListService(gameListRepository, userRepository, logger)

	ctx := context.Background()
	userID := uuid.New()
	name := "My Favorite Games"
	isPublic := true
	isDefault := false

	t.Run("should create game list successfully", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		gameListRepository.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		err := gameListService.CreateGameList(ctx, userID, name, isPublic, isDefault)

		assert.Nil(t, err)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		err := gameListService.CreateGameList(ctx, userID, name, isPublic, isDefault)

		assert.NotNil(t, err)
		assert.Equal(t, "User does not exists", err.Message)
	})

	t.Run("should return error when Find user fails", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, errors.New("database error"))

		err := gameListService.CreateGameList(ctx, userID, name, isPublic, isDefault)

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to search for user due to internal error", err.Message)
	})

	t.Run("should return error when Create game list fails", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		gameListRepository.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("database error"))

		err := gameListService.CreateGameList(ctx, userID, name, isPublic, isDefault)

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to create list due to internal error", err.Message)
	})
}

func TestGameListService_DeleteGameList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameListRepository := mocks.NewMockGameListRepository(mockCtrl)
	userRepository := mocks.NewMockUserRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameListService := service.NewGameListService(gameListRepository, userRepository, logger)

	ctx := context.Background()
	userID := uuid.New()
	gameListID := uuid.New()

	t.Run("should delete game list successfully", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{}, nil)
		gameListRepository.EXPECT().Delete(ctx, gameListID).Return(nil)

		err := gameListService.DeleteGameList(ctx, gameListID, userID)

		assert.Nil(t, err)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		err := gameListService.DeleteGameList(ctx, gameListID, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "User does not exists", err.Message)
	})

	t.Run("should return error when game list is not found", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(nil, gorm.ErrRecordNotFound)

		err := gameListService.DeleteGameList(ctx, gameListID, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "Game list does not exists", err.Message)
	})

	t.Run("should return error when Delete game list fails", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{}, nil)
		gameListRepository.EXPECT().Delete(ctx, gameListID).Return(errors.New("database error"))

		err := gameListService.DeleteGameList(ctx, gameListID, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to delete list due to internal error", err.Message)
	})
}

func TestGameListService_FindGamesFromList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameListRepository := mocks.NewMockGameListRepository(mockCtrl)
	userRepository := mocks.NewMockUserRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameListService := service.NewGameListService(gameListRepository, userRepository, logger)

	ctx := context.Background()
	gameListID := uuid.New()
	games := []*entities.Game{{ID: uuid.New()}}

	t.Run("should find games from list successfully", func(t *testing.T) {
		gameListRepository.EXPECT().FindGamesByListID(ctx, gameListID).Return(games, nil)

		foundGames, err := gameListService.FindGamesFromList(ctx, gameListID)

		assert.Nil(t, err)
		assert.Equal(t, games, foundGames)
	})

	t.Run("should return empty slice when game list is not found", func(t *testing.T) {
		gameListRepository.EXPECT().FindGamesByListID(ctx, gameListID).Return(nil, gorm.ErrRecordNotFound)

		foundGames, err := gameListService.FindGamesFromList(ctx, gameListID)

		assert.Nil(t, err)
		assert.Empty(t, foundGames)
	})

	t.Run("should return error when FindGamesByListID fails", func(t *testing.T) {
		gameListRepository.EXPECT().FindGamesByListID(ctx, gameListID).Return(nil, errors.New("database error"))

		foundGames, err := gameListService.FindGamesFromList(ctx, gameListID)

		assert.NotNil(t, err)
		assert.Nil(t, foundGames)
		assert.Equal(t, "Failed to find games in list", err.Message)
	})
}

func TestGameListService_UpdateGameList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gameListRepository := mocks.NewMockGameListRepository(mockCtrl)
	userRepository := mocks.NewMockUserRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	gameListService := service.NewGameListService(gameListRepository, userRepository, logger)

	ctx := context.Background()
	userID := uuid.New()
	gameListID := uuid.New()
	name := "My Updated Favorite Games"
	isPublic := false

	t.Run("should update game list successfully", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		gameListRepository.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		err := gameListService.UpdateGameList(ctx, userID, gameListID, name, isPublic)

		assert.Nil(t, err)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		err := gameListService.UpdateGameList(ctx, userID, gameListID, name, isPublic)

		assert.NotNil(t, err)
		assert.Equal(t, "User does not exist", err.Message)
	})

	t.Run("should return error when Find user fails", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, errors.New("database error"))

		err := gameListService.UpdateGameList(ctx, userID, gameListID, name, isPublic)

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to search for user due to internal error", err.Message)
	})

	t.Run("should return error when Update game list fails", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		gameListRepository.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("database error"))

		err := gameListService.UpdateGameList(ctx, userID, gameListID, name, isPublic)

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to create list due to internal error", err.Message)
	})
}
