package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	userIDKey     = "user_id"
	gameIDKey     = "game_id"
	gameListIDKey = "game_list_id"
)

func GetUserID(c echo.Context) uuid.UUID {
	return c.Get(string(userIDKey)).(uuid.UUID)
}

func GetGameID(c echo.Context) uuid.UUID {
	return c.Get(string(gameIDKey)).(uuid.UUID)
}

func GetGameListID(c echo.Context) uuid.UUID {
	return c.Get(string(gameListIDKey)).(uuid.UUID)
}
