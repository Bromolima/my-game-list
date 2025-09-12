package middlewares

import (
	"log/slog"

	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/token"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	jwtService token.JwtService
}

func NewAuthMiddleware(jwtService token.JwtService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) TokenAuthentificationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	log := slog.With(
		slog.String("auth", "middleware"),
		slog.String("func", "TokenAuthenticationMiddleware"),
	)

	return func(ectx echo.Context) error {
		if err := m.jwtService.VerifyToken(ectx); err != nil {
			log.Warn("failed to validate token", slog.String("error", err.Error()))
			restErr := resterr.NewUnauthorizedError("failed validating token")
			return ectx.JSON(restErr.Code, restErr)
		}

		return next(ectx)
	}
}
