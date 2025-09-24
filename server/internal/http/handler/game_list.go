package handler

import (
	"log/slog"
	"net/http"

	"github.com/Bromolima/my-game-list/internal/factory"
	"github.com/Bromolima/my-game-list/internal/http/dto"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/service"
	"github.com/Bromolima/my-game-list/internal/token"
	"github.com/Bromolima/my-game-list/internal/validation"
	"github.com/labstack/echo/v4"
)

type GameListHandler struct {
	gameListService service.GameListService
	logger          *slog.Logger
	jwtService      token.JwtService
}

func NewGameListHandler(gameListService service.GameListService, logger *slog.Logger) *GameListHandler {
	return &GameListHandler{
		gameListService: gameListService,
		logger:          logger.With(slog.String("handler", "gameList")),
	}
}

func (h *GameListHandler) CreateGameList(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "CreateGameList"))

	var createRequest dto.GameListCreateRequest
	if err := ectx.Bind(&createRequest); err != nil {
		log.Warn("Failed to bind request payload", slog.String("error", err.Error()))
		restErr := resterr.NewBadRequestError("An error occurred while binding the request payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(createRequest); err != nil {
		log.Warn("Request payload validation failed", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	userClaims, err := h.jwtService.ExtractToken(ectx)
	if err != nil {
		log.Warn("Failed to extract token from context")
		restErr := resterr.NewUnauthorizedError("Failed to extract token")
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.gameListService.CreateGameList(
		ectx.Request().Context(),
		userClaims.ID,
		createRequest.Name,
		createRequest.IsPublic,
		false,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game list created successfully")
	return ectx.NoContent(http.StatusCreated)
}

func (h *GameListHandler) FindGamesFromList(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "FindGamesFromList"))

	gameListID := GetGameListID(ectx)
	gameList, restErr := h.gameListService.FindGamesFromList(ectx.Request().Context(), gameListID)
	if restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	var gamesResponse []*dto.GameResponse
	for _, game := range gameList {
		gamesResponse = append(gamesResponse, factory.NewResponseFromGame(game))
	}

	log.Info("Game list found successfully")
	return ectx.JSON(http.StatusOK, gamesResponse)
}

func (h *GameListHandler) UpdateGameList(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "UpdateGameList"))

	var updateRequest dto.GameListUpdateRequest
	if err := ectx.Bind(&updateRequest); err != nil {
		log.Warn("Failed to bind request payload", slog.String("error", err.Error()))
		restErr := resterr.NewBadRequestError("An error occurred while binding the request payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(updateRequest); err != nil {
		log.Warn("Request payload validation failed", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	userClaims, err := h.jwtService.ExtractToken(ectx)
	if err != nil {
		log.Warn("Failed to extract token from context")
		restErr := resterr.NewUnauthorizedError("Failed to extract token")
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.gameListService.UpdateGameList(
		ectx.Request().Context(),
		userClaims.ID,
		GetGameListID(ectx),
		*updateRequest.Name,
		*updateRequest.IsPublic,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game list updated successfully")
	return ectx.NoContent(http.StatusOK)
}

func (h *GameListHandler) DeleteGameList(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "DeleteGameList"))

	gameListID := GetGameListID(ectx)
	userClaims, err := h.jwtService.ExtractToken(ectx)
	if err != nil {
		log.Warn("Failed to extract token from context")
		restErr := resterr.NewUnauthorizedError("Failed to extract token")
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.gameListService.DeleteGameList(ectx.Request().Context(), gameListID, userClaims.ID); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game list deleted successfully")
	return ectx.NoContent(http.StatusOK)
}
