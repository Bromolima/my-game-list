package routes

import (
	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/http/handler"
	"github.com/Bromolima/my-game-list/internal/http/middlewares"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func SetupRoutes(e *echo.Echo, c *dig.Container) error {
	if err := setupUserRoutes(e, c); err != nil {
		return err
	}

	if err := setupAuthRoutes(e, c); err != nil {
		return err
	}

	if err := setupGamesRoutes(e, c); err != nil {
		return err
	}

	if err := setupGameListRoutes(e, c); err != nil {
		return err
	}

	if err := setupListItemRoutes(e, c); err != nil {
		return err
	}

	return nil
}

func setupAuthRoutes(e *echo.Echo, c *dig.Container) error {
	return c.Invoke(func(h *handler.UserHandler, m *middlewares.AuthMiddleware) {
		g := e.Group("/auth")

		g.POST("/register", h.RegisterUser)
		g.POST("/login", h.Login)
	})
}

func setupUserRoutes(e *echo.Echo, c *dig.Container) error {
	return c.Invoke(func(h *handler.UserHandler, m *middlewares.AuthMiddleware) {
		g := e.Group("/users")

		g.GET("/", h.SearchUsers)
		g.PUT("/:id", h.UpdateUser, m.AuthMiddleware)
		g.DELETE("/:id", h.DeleteUser, m.AuthMiddleware)
	})
}

func setupGamesRoutes(e *echo.Echo, c *dig.Container) error {
	return c.Invoke(func(h *handler.GameHandler, m *middlewares.AuthMiddleware) {
		g := e.Group("/games")

		g.POST("", h.CreateGame, m.RequireAccess(entities.CreateAcess))
		g.GET("", h.SearchGames, m.RequireAccess(entities.ReadAccess))
		g.PUT("/:id", h.UpdateGame, m.RequireAccess(entities.UpdateAccess))
		g.DELETE("/id", h.DeleteGame, m.RequireAccess(entities.DeleteAcess))
	})
}

func setupGameListRoutes(e *echo.Echo, c *dig.Container) error {
	return c.Invoke(func(h *handler.GameListHandler, m *middlewares.AuthMiddleware) {
		g := e.Group("/list")

		g.POST("", h.CreateGameList)
		g.GET("/:id", h.FindGamesFromList)
		g.PUT("", h.UpdateGameList)
		g.DELETE("", h.DeleteGameList)
	})
}

func setupListItemRoutes(e *echo.Echo, c *dig.Container) error {
	return c.Invoke(func(h *handler.ListItemHandler, m *middlewares.AuthMiddleware) {
		e.POST("", h.AddGameToList, m.AuthMiddleware)
		e.PUT("", h.UpdateGameFromList, m.AuthMiddleware)
		e.DELETE("", h.DeleteGameFromList, m.AuthMiddleware)
	})
}
