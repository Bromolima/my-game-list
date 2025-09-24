package handler

import (
	"log/slog"
	"net/http"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/factory"
	"github.com/Bromolima/my-game-list/internal/http/dto"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/service"
	"github.com/Bromolima/my-game-list/internal/validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GameHandler struct {
	gameService service.GameService
	logger      *slog.Logger
}

func NewGameHandler(gameService service.GameService, logger *slog.Logger) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		logger:      logger.With(slog.String("handler", "game")),
	}
}

func (h *GameHandler) CreateGame(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "CreateGame"))

	var createRequest dto.GameCreateRequest
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

	if restErr := h.gameService.CreateGame(
		ectx.Request().Context(),
		createRequest.Name,
		createRequest.Genre,
		createRequest.Developer,
		createRequest.Description,
		createRequest.ImageURL,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game created successfully")
	return ectx.NoContent(http.StatusCreated)
}

func (h *GameHandler) SearchGames(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "SearchGames"))

	var searchRequest dto.GamesSearchRequest
	if err := ectx.Bind(&searchRequest); err != nil {
		log.Warn("Failed to bind request payload")
		restErr := resterr.NewBadRequestError("An error occurred while binding the request payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	page, restErr := h.gameService.SearchGames(
		ectx.Request().Context(),
		factory.NewPage[entities.Game](searchRequest.Page, searchRequest.Limit),
		searchRequest.Name,
	)
	if restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	pageResonse := factory.NewReponseFromPage(page, factory.NewResponseFromGame)

	log.Info("Games searched successfully")
	return ectx.JSON(http.StatusOK, pageResonse)
}

func (h *GameHandler) UpdateGame(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "UpdateGame"))

	id := ectx.Param("id")
	var updateRequest dto.GameUpdateRequest
	if err := ectx.Bind(&updateRequest); err != nil {
		log.Warn("Failed to bind request payload", slog.String("error", err.Error()))
		restErr := resterr.NewBadRequestError("An error occurred while binding the request payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(&updateRequest); err != nil {
		log.Warn("Request payload validation failed", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	gameID, err := uuid.Parse(id)
	if err != nil {
		log.Error("Failed to parse game ID from path parameter", slog.String("error", err.Error()))
		restErr := resterr.NewBadRequestError("An error ocurried while parsing the id")
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.gameService.UpdateGame(
		ectx.Request().Context(),
		gameID,
		*updateRequest.Name,
		*updateRequest.Genre,
		*updateRequest.Developer,
		*updateRequest.Description,
		*updateRequest.ImageURL,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game updated successfully")
	return ectx.NoContent(http.StatusOK)
}

func (h *GameHandler) DeleteGame(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "DeleteGame"))

	id := ectx.Param("id")
	if restErr := h.gameService.DeleteGame(ectx.Request().Context(), id); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game deleted successfully")
	return ectx.NoContent(http.StatusOK)
}
