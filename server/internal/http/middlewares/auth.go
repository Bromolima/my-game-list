package middlewares

import (
	"log/slog"

	"github.com/Bromolima/my-game-list/internal/entities"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/repository"
	"github.com/Bromolima/my-game-list/internal/token"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	jwtService     token.JwtService
	roleRepository repository.RoleRepository
	loggerr        *slog.Logger
}

func NewAuthMiddleware(jwtService token.JwtService, roleRepository repository.RoleRepository, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:     jwtService,
		roleRepository: roleRepository,
		loggerr:        logger.With(slog.String("middleware", "auth")),
	}
}

func (m *AuthMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	log := m.loggerr.With(slog.String("func", "AuthMiddleware"))

	return func(ectx echo.Context) error {
		if err := m.jwtService.ValidateToken(ectx); err != nil {
			log.Warn("An error occurred while validating the token")
			restErr := resterr.NewUnauthorizedError("An error occurred while validating the token")
			return ectx.JSON(restErr.Code, restErr)
		}

		return next(ectx)
	}
}

func (m *AuthMiddleware) RequireAccess(access entities.AccessType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log := m.loggerr.With(slog.String("func", "RequireAccess"))

			if err := m.jwtService.ValidateToken(c); err != nil {
				log.Warn("Invalid token")
				restErr := resterr.NewUnauthorizedError("Invalid token")
				return c.JSON(restErr.Code, restErr)
			}

			userClaims, err := m.jwtService.ExtractToken(c)
			if err != nil {
				log.Warn("Invalid token claims")
				restErr := resterr.NewUnauthorizedError("Invalid token")
				return c.JSON(restErr.Code, restErr)
			}

			err, permitted := m.roleRepository.HasAccess(c.Request().Context(), userClaims.ID, access)
			if err != nil {
				log.Error("Error checking access", "error", err.Error())
				restErr := resterr.NewInternalServerErr("Failed to validate access")
				return c.JSON(restErr.Code, restErr)
			}

			if !permitted {
				log.Warn("User forbidden to access this feature")
				restErr := resterr.NewForbiddenError("User is forbidden to access this feature")
				return c.JSON(restErr.Code, restErr)
			}

			return next(c)
		}
	}
}
