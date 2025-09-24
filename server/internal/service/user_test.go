package service_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/security"
	"github.com/Bromolima/my-game-list/internal/service"
	"github.com/Bromolima/my-game-list/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestUserService_RegisterUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepository := mocks.NewMockUserRepository(mockCtrl)
	tokenService := mocks.NewMockJwtService(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userService := service.NewUserService(userRepository, tokenService, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	username := "testuser"
	avatarURL := "http://example.com/avatar.png"

	t.Run("should register user successfully", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)
		userRepository.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		err := userService.RegisterUser(ctx, email, password, username, avatarURL)

		assert.Nil(t, err)
	})

	t.Run("should return error when user already exists", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(&entities.User{}, nil)

		err := userService.RegisterUser(ctx, email, password, username, avatarURL)

		assert.NotNil(t, err)
		assert.Equal(t, "The provided email is already registered", err.Message)
	})

	t.Run("should return error when FindByEmail fails", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(nil, errors.New("database error"))

		err := userService.RegisterUser(ctx, email, password, username, avatarURL)

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while finding the user by email", err.Message)
	})

	t.Run("should return error when Create fails", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)
		userRepository.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("database error"))

		err := userService.RegisterUser(ctx, email, password, username, avatarURL)

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while creating the user", err.Message)
	})
}

func TestUserService_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepository := mocks.NewMockUserRepository(mockCtrl)
	tokenService := mocks.NewMockJwtService(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userService := service.NewUserService(userRepository, tokenService, logger)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := security.HashPassword(password)
	user := &entities.User{
		Email:    email,
		Password: hashedPassword,
	}

	t.Run("should login successfully", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(user, nil)
		tokenService.EXPECT().GenerateToken(user).Return("token", nil)

		token, err := userService.Login(ctx, email, password)

		assert.Nil(t, err)
		assert.Equal(t, "token", token)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)

		token, err := userService.Login(ctx, email, password)

		assert.NotNil(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "The requested user was not found", err.Message)
	})

	t.Run("should return error when password is wrong", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(user, nil)

		token, err := userService.Login(ctx, email, "wrongpassword")

		assert.NotNil(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "Invalid credentials provided", err.Message)
	})

	t.Run("should return error when FindByEmail fails", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(nil, errors.New("database error"))

		token, err := userService.Login(ctx, email, password)

		assert.NotNil(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "An error occurred while finding the user", err.Message)
	})

	t.Run("should return error when GenerateToken fails", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(user, nil)
		tokenService.EXPECT().GenerateToken(user).Return("", errors.New("token error"))

		token, err := userService.Login(ctx, email, password)

		assert.NotNil(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "An error occurred while generating the token", err.Message)
	})
}

func TestUserService_FindUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepository := mocks.NewMockUserRepository(mockCtrl)
	tokenService := mocks.NewMockJwtService(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userService := service.NewUserService(userRepository, tokenService, logger)

	ctx := context.Background()
	userID := uuid.New()
	user := &entities.User{ID: userID}

	t.Run("should find user successfully", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(user, nil)

		foundUser, err := userService.FindUser(ctx, userID.String())

		assert.Nil(t, err)
		assert.Equal(t, user, foundUser)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		foundUser, err := userService.FindUser(ctx, userID.String())

		assert.NotNil(t, err)
		assert.Nil(t, foundUser)
		assert.Equal(t, "The requested user was not found", err.Message)
	})

	t.Run("should return error when Find fails", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, errors.New("database error"))

		foundUser, err := userService.FindUser(ctx, userID.String())

		assert.NotNil(t, err)
		assert.Nil(t, foundUser)
		assert.Equal(t, "An error occurred while finding the user", err.Message)
	})
}

func TestUserService_SearchUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepository := mocks.NewMockUserRepository(mockCtrl)
	tokenService := mocks.NewMockJwtService(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userService := service.NewUserService(userRepository, tokenService, logger)

	ctx := context.Background()
	page := &entities.Page[entities.User]{}
	query := "test"

	t.Run("should search users successfully", func(t *testing.T) {
		userRepository.EXPECT().Search(ctx, page, query).Return(page, nil)

		resultPage, err := userService.SearchUsers(ctx, page, query)

		assert.Nil(t, err)
		assert.Equal(t, page, resultPage)
	})

	t.Run("should return error when no users are found", func(t *testing.T) {
		userRepository.EXPECT().Search(ctx, page, query).Return(nil, gorm.ErrRecordNotFound)

		resultPage, err := userService.SearchUsers(ctx, page, query)

		assert.NotNil(t, err)
		assert.Nil(t, resultPage)
		assert.Equal(t, "No users were found for the given query", err.Message)
	})

	t.Run("should return error when Search fails", func(t *testing.T) {
		userRepository.EXPECT().Search(ctx, page, query).Return(nil, errors.New("database error"))

		resultPage, err := userService.SearchUsers(ctx, page, query)

		assert.NotNil(t, err)
		assert.Nil(t, resultPage)
		assert.Equal(t, "An error occurred while searching for users", err.Message)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepository := mocks.NewMockUserRepository(mockCtrl)
	tokenService := mocks.NewMockJwtService(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userService := service.NewUserService(userRepository, tokenService, logger)

	ctx := context.Background()
	userID := uuid.New()
	email := "test@example.com"
	password := "password123"
	username := "testuser"
	avatarURL := "http://example.com/avatar.png"

	t.Run("should update user successfully", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(&entities.User{}, nil)
		userRepository.EXPECT().Update(ctx, gomock.Any()).Return(nil)

		err := userService.UpdateUser(ctx, userID, email, password, username, avatarURL)

		assert.Nil(t, err)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(nil, gorm.ErrRecordNotFound)

		err := userService.UpdateUser(ctx, userID, email, password, username, avatarURL)

		assert.NotNil(t, err)
		assert.Equal(t, "The requested user was not found", err.Message)
	})

	t.Run("should return error when FindByEmail fails", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(nil, errors.New("database error"))

		err := userService.UpdateUser(ctx, userID, email, password, username, avatarURL)

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while finding the user", err.Message)
	})

	t.Run("should return error when Update fails", func(t *testing.T) {
		userRepository.EXPECT().FindByEmail(ctx, email).Return(&entities.User{}, nil)
		userRepository.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("database error"))

		err := userService.UpdateUser(ctx, userID, email, password, username, avatarURL)

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while updating the user", err.Message)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepository := mocks.NewMockUserRepository(mockCtrl)
	tokenService := mocks.NewMockJwtService(mockCtrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userService := service.NewUserService(userRepository, tokenService, logger)

	ctx := context.Background()
	userID := uuid.New()

	t.Run("should delete user successfully", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		userRepository.EXPECT().Delete(ctx, userID).Return(nil)

		err := userService.DeleteUser(ctx, userID.String())

		assert.Nil(t, err)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		err := userService.DeleteUser(ctx, userID.String())

		assert.NotNil(t, err)
		assert.Equal(t, "The requested user was not found", err.Message)
	})

	t.Run("should return error when Find fails", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(nil, errors.New("database error"))

		err := userService.DeleteUser(ctx, userID.String())

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while finding the user", err.Message)
	})

	t.Run("should return error when Delete fails", func(t *testing.T) {
		userRepository.EXPECT().Find(ctx, userID).Return(&entities.User{}, nil)
		userRepository.EXPECT().Delete(ctx, userID).Return(errors.New("database error"))

		err := userService.DeleteUser(ctx, userID.String())

		assert.NotNil(t, err)
		assert.Equal(t, "An error occurred while deleting the user", err.Message)
	})
}
