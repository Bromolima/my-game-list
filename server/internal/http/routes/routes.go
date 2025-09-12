package routes

import (
	"github.com/Bromolima/my-game-list/internal/http/handler"
	"github.com/Bromolima/my-game-list/internal/http/middlewares"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func SetupRoutes(e *echo.Echo, c *dig.Container) error {
	if err := c.Invoke(func(h *handler.UserHandler, m *middlewares.AuthMiddleware) {
		e.POST("/users", h.RegisterUser)
		e.GET("/users/:id", h.FindUser, m.TokenAuthentificationMiddleware)
		e.PUT("/users/:id", h.UpdateUser, m.TokenAuthentificationMiddleware)
		e.DELETE("/users/:id", h.DeleteUser, m.TokenAuthentificationMiddleware)

		e.POST("/login", h.Login)
	}); err != nil {
		return err
	}

	return nil
}
