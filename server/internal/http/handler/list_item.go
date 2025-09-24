package handler

import (
	"log/slog"
	"net/http"

	"github.com/Bromolima/my-game-list/internal/http/dto"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/service"
	"github.com/Bromolima/my-game-list/internal/token"
	"github.com/Bromolima/my-game-list/internal/validation"
	"github.com/labstack/echo/v4"
)

type ListItemHandler struct {
	listItemService service.ListItemService
	jwtService      token.JwtService
	logger          *slog.Logger
}

func NewListItemHandler(listItemService service.ListItemService, jwtService token.JwtService, logger *slog.Logger) *ListItemHandler {
	return &ListItemHandler{
		listItemService: listItemService,
		jwtService:      jwtService,
		logger:          logger.With(slog.String("handler", "listItem")),
	}
}

func (h *ListItemHandler) AddGameToList(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "AddGameToList"))

	userClaims, err := h.jwtService.ExtractToken(ectx)
	if err != nil {
		log.Warn("Failed to extract user claims from token")
		restErr := resterr.NewForbiddenError("Failed to get claims")
		return ectx.JSON(restErr.Code, restErr)
	}

	var addRequest dto.ListItemAddRequest
	if err := ectx.Bind(&addRequest); err != nil {
		log.Warn("Failed to bind request payload")
		restErr := resterr.NewBadRequestError("Failed to unmarshal request")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(addRequest); err != nil {
		log.Warn("Request payload validation failed", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.listItemService.AddGameToList(
		ectx.Request().Context(),
		userClaims.ID,
		addRequest.ListID,
		addRequest.GameID,
		addRequest.Status,
		addRequest.Rating,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game added to list successfully")
	return ectx.NoContent(http.StatusCreated)
}

func (h *ListItemHandler) UpdateGameFromList(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "UpdateGameFromList"))

	userClaims, err := h.jwtService.ExtractToken(ectx)
	if err != nil {
		log.Warn("Failed to extract user claims from token")
		restErr := resterr.NewForbiddenError("Failed to get claims")
		return ectx.JSON(restErr.Code, restErr)
	}

	var updateRequest dto.ListItemUpdateRequest
	if err := ectx.Bind(&updateRequest); err != nil {
		log.Warn("Failed to bind request payload")
		restErr := resterr.NewBadRequestError("Failed to unmarshal request")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(updateRequest); err != nil {
		log.Warn("Request payload validation failed", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.listItemService.UpdateGameFromList(
		ectx.Request().Context(),
		updateRequest.GameID,
		updateRequest.ListID,
		userClaims.ID,
		*updateRequest.Rating,
		*updateRequest.Status,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game updated in list successfully")
	return ectx.NoContent(http.StatusOK)
}

func (h *ListItemHandler) DeleteGameFromList(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "DeleteGameFromList"))

	userClaims, err := h.jwtService.ExtractToken(ectx)
	if err != nil {
		log.Warn("Failed to extract user claims from token")
		restErr := resterr.NewForbiddenError("Failed to get claims")
		return ectx.JSON(restErr.Code, restErr)
	}

	var deleteRequest dto.ListItemDeleteRequest
	if err := ectx.Bind(&deleteRequest); err != nil {
		log.Warn("Failed to bind request payload")
		restErr := resterr.NewBadRequestError("Failed to unmarshal request")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(deleteRequest); err != nil {
		log.Warn("Request payload validation failed", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.listItemService.DeleteGameFromList(
		ectx.Request().Context(),
		deleteRequest.GameID,
		deleteRequest.ListID,
		userClaims.ID,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("Game deleted from list successfully")
	return ectx.NoContent(http.StatusOK)
}
