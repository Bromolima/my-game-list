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

func TestListItemService_AddGameToList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	listItemRepository := mocks.NewMockListItemRepository(mockCtrl)
	gameRepository := mocks.NewMockGameRepository(mockCtrl)
	gameListRepository := mocks.NewMockGameListRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	listItemService := service.NewListItemService(listItemRepository, gameRepository, gameListRepository, logger)

	ctx := context.Background()
	userID := uuid.New()
	gameID := uuid.New()
	gameListID := uuid.New()
	status := "Playing"
	rating := float32(4.5)

	t.Run("should add game to list successfully", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: userID}, nil)
		listItemRepository.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		err := listItemService.AddGameToList(ctx, userID, gameID, gameListID, status, rating)

		assert.Nil(t, err)
	})

	t.Run("should return error when game is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, gorm.ErrRecordNotFound)

		err := listItemService.AddGameToList(ctx, userID, gameID, gameListID, status, rating)

		assert.NotNil(t, err)
		assert.Equal(t, "The specified game was not found", err.Message)
	})

	t.Run("should return error when game list is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(nil, gorm.ErrRecordNotFound)

		err := listItemService.AddGameToList(ctx, userID, gameID, gameListID, status, rating)

		assert.NotNil(t, err)
		assert.Equal(t, "The specified game list was not found", err.Message)
	})

	t.Run("should return error when user is not authorized", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: uuid.New()}, nil)

		err := listItemService.AddGameToList(ctx, userID, gameID, gameListID, status, rating)

		assert.NotNil(t, err)
		assert.Equal(t, "You do not have permission to modify this list", err.Message)
	})

	t.Run("should return error when Create fails", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: userID}, nil)
		listItemRepository.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("database error"))

		err := listItemService.AddGameToList(ctx, userID, gameID, gameListID, status, rating)

		assert.NotNil(t, err)
		assert.Equal(t, "Could not add the game to the list", err.Message)
	})
}

func TestListItemService_DeleteGameFromList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	listItemRepository := mocks.NewMockListItemRepository(mockCtrl)
	gameRepository := mocks.NewMockGameRepository(mockCtrl)
	gameListRepository := mocks.NewMockGameListRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	listItemService := service.NewListItemService(listItemRepository, gameRepository, gameListRepository, logger)

	ctx := context.Background()
	userID := uuid.New()
	gameID := uuid.New()
	gameListID := uuid.New()

	t.Run("should delete game from list successfully", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: userID}, nil)
		listItemRepository.EXPECT().Delete(ctx, gameID, gameListID).Return(nil)

		err := listItemService.DeleteGameFromList(ctx, gameID, gameListID, userID)

		assert.Nil(t, err)
	})

	t.Run("should return error when game is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, gorm.ErrRecordNotFound)

		err := listItemService.DeleteGameFromList(ctx, gameID, gameListID, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "The specified game was not found", err.Message)
	})

	t.Run("should return error when game list is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(nil, gorm.ErrRecordNotFound)

		err := listItemService.DeleteGameFromList(ctx, gameID, gameListID, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "The specified game list was not found", err.Message)
	})

	t.Run("should return error when user is not authorized", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: uuid.New()}, nil)

		err := listItemService.DeleteGameFromList(ctx, gameID, gameListID, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "You do not have permission to modify this list", err.Message)
	})

	t.Run("should return error when Delete fails", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: userID}, nil)
		listItemRepository.EXPECT().Delete(ctx, gameID, gameListID).Return(errors.New("database error"))

		err := listItemService.DeleteGameFromList(ctx, gameID, gameListID, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "Could not add the game to the list", err.Message)
	})
}

func TestListItemService_UpdateGameFromList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	listItemRepository := mocks.NewMockListItemRepository(mockCtrl)
	gameRepository := mocks.NewMockGameRepository(mockCtrl)
	gameListRepository := mocks.NewMockGameListRepository(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	listItemService := service.NewListItemService(listItemRepository, gameRepository, gameListRepository, logger)

	ctx := context.Background()
	userID := uuid.New()
	gameID := uuid.New()
	gameListID := uuid.New()
	status := "Finished"
	rating := float32(5)

	t.Run("should update game from list successfully", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: userID}, nil)
		listItemRepository.EXPECT().Update(ctx, gameID, gameListID, rating, status).Return(nil)

		err := listItemService.UpdateGameFromList(ctx, gameID, gameListID, userID, rating, status)

		assert.Nil(t, err)
	})

	t.Run("should return error when game is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(nil, gorm.ErrRecordNotFound)

		err := listItemService.UpdateGameFromList(ctx, gameID, gameListID, userID, rating, status)

		assert.NotNil(t, err)
		assert.Equal(t, "The specified game was not found", err.Message)
	})

	t.Run("should return error when game list is not found", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(nil, gorm.ErrRecordNotFound)

		err := listItemService.UpdateGameFromList(ctx, gameID, gameListID, userID, rating, status)

		assert.NotNil(t, err)
		assert.Equal(t, "The specified game list was not found", err.Message)
	})

	t.Run("should return error when user is not authorized", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: uuid.New()}, nil)

		err := listItemService.UpdateGameFromList(ctx, gameID, gameListID, userID, rating, status)

		assert.NotNil(t, err)
		assert.Equal(t, "You do not have permission to modify this list", err.Message)
	})

	t.Run("should return error when Update fails", func(t *testing.T) {
		gameRepository.EXPECT().Find(ctx, gameID).Return(&entities.Game{}, nil)
		gameListRepository.EXPECT().Find(ctx, gameListID).Return(&entities.GameList{UserID: userID}, nil)
		listItemRepository.EXPECT().Update(ctx, gameID, gameListID, rating, status).Return(errors.New("database error"))

		err := listItemService.UpdateGameFromList(ctx, gameID, gameListID, userID, rating, status)

		assert.NotNil(t, err)
		assert.Equal(t, "Could not add the game to the list", err.Message)
	})
}
