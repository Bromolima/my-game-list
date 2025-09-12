package handler

import (
	"log/slog"
	"net/http"

	"github.com/Bromolima/my-game-list/internal/http/cookie"
	"github.com/Bromolima/my-game-list/internal/http/dto"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/models"
	"github.com/Bromolima/my-game-list/internal/validation"
	"github.com/labstack/echo/v4"
)

func (h *UserHandler) Login(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "Login"))

	var loginPayload dto.UserLoginPayload
	if err := ectx.Bind(&loginPayload); err != nil {
		log.Warn("failed to bind request payload", slog.String("error", err.Error()))
		restErr := resterr.NewBadRequestError("Failed to bind payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(loginPayload); err != nil {
		log.Warn("failed to validate request payload", slog.String("error", err.Error()))
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	user, token, err := h.service.Login(ectx.Request().Context(), &models.User{
		Email:    loginPayload.Email,
		Password: loginPayload.Password,
	})

	if err != nil {
		return ectx.JSON(err.Code, err)
	}

	cookie.SetCookie(ectx, token)
	log.Info("user logged in successfully", slog.String("username", user.Username))
	return ectx.JSON(http.StatusOK, &dto.UserResponse{
		Username: user.Username,
	})
}
