package handler

import (
	"log/slog"
	"net/http"

	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/factory"
	"github.com/Bromolima/my-game-list/internal/http/cookie"
	"github.com/Bromolima/my-game-list/internal/http/dto"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/service"
	"github.com/Bromolima/my-game-list/internal/token"
	validation "github.com/Bromolima/my-game-list/internal/validation"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService     service.UserService
	gameListService service.GameListService
	jwtService      token.JwtService
	logger          *slog.Logger
}

func NewUserHandler(userService service.UserService, gameListService service.GameListService, jwtService token.JwtService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService:     userService,
		gameListService: gameListService,
		logger:          logger.With(slog.String("handler", "user")),
		jwtService:      jwtService,
	}
}

func (h *UserHandler) RegisterUser(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "RegisterUser"))

	var registerRequest dto.UserRegisterRequest
	if err := ectx.Bind(&registerRequest); err != nil {
		log.Warn("Failed to bind request payload")
		restErr := resterr.NewBadRequestError("invalid payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(registerRequest); err != nil {
		log.Warn("Request payload validation failed")
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.userService.RegisterUser(
		ectx.Request().Context(),
		registerRequest.Email,
		registerRequest.Password,
		registerRequest.Username,
		registerRequest.AvatarURL,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("User registered successfully")
	return ectx.NoContent(http.StatusCreated)
}

func (h *UserHandler) Login(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "Login"))

	var loginPayload dto.UserLoginRequest
	if err := ectx.Bind(&loginPayload); err != nil {
		log.Warn("Failed to bind request payload")
		restErr := resterr.NewBadRequestError("An error occurred while binding the request payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	if err := ectx.Validate(loginPayload); err != nil {
		log.Warn("Request payload validation failed")
		restErr := validation.ValidateUserError(err)
		return ectx.JSON(restErr.Code, restErr)
	}

	token, restErr := h.userService.Login(ectx.Request().Context(), loginPayload.Email, loginPayload.Password)

	if restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	cookie.SetCookie(ectx, token)
	log.Info("User logged in successfully")
	return ectx.NoContent(http.StatusOK)
}

func (h *UserHandler) SearchUsers(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "SearchUsers"))

	var searchRequest dto.UserSearchRequest
	if err := ectx.Bind(&searchRequest); err != nil {
		log.Warn("Failed to bind request payload")
		restErr := resterr.NewBadRequestError("An error occurred while binding the request payload")
		return ectx.JSON(restErr.Code, restErr)
	}

	page, restErr := h.userService.SearchUsers(
		ectx.Request().Context(),
		factory.NewPage[entities.User](searchRequest.Page, searchRequest.Limit),
		searchRequest.Name,
	)
	if restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	pageResonse := factory.NewReponseFromPage(page, factory.NewResponseFromUser)

	log.Info("Users searched successfully")
	return ectx.JSON(http.StatusOK, pageResonse)
}

func (h *UserHandler) UpdateUser(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "UpdateUser"))

	userClaims, err := h.jwtService.ExtractToken(ectx)
	if err != nil {
		log.Warn("Failed to extract token from context")
		restErr := resterr.NewUnauthorizedError("Failed to extract token")
		return ectx.JSON(restErr.Code, restErr)
	}

	var updateRequest dto.UserUpdateRequest
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

	if restErr := h.userService.UpdateUser(
		ectx.Request().Context(),
		userClaims.ID,
		*updateRequest.Email,
		*updateRequest.Password,
		*updateRequest.Username,
		*updateRequest.AvatarURL,
	); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("User updated successfully")
	return ectx.NoContent(http.StatusNoContent)
}

func (h *UserHandler) DeleteUser(ectx echo.Context) error {
	log := h.logger.With(slog.String("func", "DeleteUser"))

	userClaims, err := h.jwtService.ExtractToken(ectx)
	if err != nil {
		log.Warn("Failed to extract token from context")
		restErr := resterr.NewUnauthorizedError("Failed to extract token")
		return ectx.JSON(restErr.Code, restErr)
	}

	if restErr := h.userService.DeleteUser(ectx.Request().Context(), userClaims.ID.String()); restErr != nil {
		return ectx.JSON(restErr.Code, restErr)
	}

	log.Info("User deleted successfully")
	return ectx.NoContent(http.StatusNoContent)
}
