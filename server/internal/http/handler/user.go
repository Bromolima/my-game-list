package handler

import (
	"log/slog"
	"net/http"

	"github.com/Bromolima/my-game-list/internal/http/dto"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/models"
	"github.com/Bromolima/my-game-list/internal/service"
	validation "github.com/Bromolima/my-game-list/internal/validation"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

func NewUserHandler(service service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger.With(slog.String("handler", "user")),
	}
}

func (h *UserHandler) RegisterUser(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "RegisterUser"))

	var payload dto.UserCreatePayload
	if err := ectx.Bind(&payload); err != nil {
		log.Warn("failed to bind request payload", slog.String("error", err.Error()))
		restErr := resterr.NewBadRequestError("Failed to bind payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(payload); err != nil {
		log.Warn("failed to validate request payload", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := h.service.CreateUser(ectx.Request().Context(), &models.User{
		Email:    payload.Email,
		Password: payload.Password,
		Username: payload.Username,
	}); err != nil {
		return ectx.JSON(err.Code, err)
	}

	log.Info("user registered successfully", slog.String("username", payload.Username))
	return ectx.JSON(http.StatusCreated, &dto.UserResponse{
		Username: payload.Username,
	})
}

func (h *UserHandler) FindUser(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "FindUser"))

	id := ectx.Param("id")
	user, err := h.service.FindUser(ectx.Request().Context(), id)
	if err != nil {
		return ectx.JSON(err.Code, err)
	}

	log.Info("user found successfully", slog.String("username", user.Username))
	return ectx.JSON(http.StatusFound, &dto.UserResponse{
		Username: user.Username,
	})
}

func (h *UserHandler) UpdateUser(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "UpdateUser"))

	id := ectx.Param("id")
	var payload dto.UserCreatePayload
	if err := ectx.Bind(&payload); err != nil {
		log.Warn("failed to bind request payload", slog.String("error", err.Error()))
		restErr := resterr.NewBadRequestError("Failed to bind payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(payload); err != nil {
		log.Warn("failed to validate request payload", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := h.service.UpdateUser(ectx.Request().Context(), &models.User{
		ID:       id,
		Email:    payload.Email,
		Password: payload.Password,
		Username: payload.Username,
	}); err != nil {
		return ectx.JSON(err.Code, err)
	}

	log.Info("user updated successfully", slog.String("username", payload.Username))
	return ectx.JSON(http.StatusOK, &dto.UserResponse{
		Username: payload.Username,
	})
}

func (h *UserHandler) DeleteUser(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "DeleteUser"))

	id := ectx.Param("id")
	if err := h.service.DeleteUser(ectx.Request().Context(), id); err != nil {
		return ectx.JSON(err.Code, err)
	}

	log.Info("user deleted successfully", slog.String("id", id))
	return ectx.NoContent(http.StatusNoContent)
}
