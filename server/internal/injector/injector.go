package injector

import (
	"github.com/Bromolima/my-game-list/internal/http/handler"
	"github.com/Bromolima/my-game-list/internal/http/middlewares"
	"github.com/Bromolima/my-game-list/internal/repository"
	"github.com/Bromolima/my-game-list/internal/service"
	"github.com/Bromolima/my-game-list/internal/token"
	"github.com/Bromolima/my-game-list/logger"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func SetupDependecies(c *dig.Container, db *gorm.DB) {
	c.Provide(func() *gorm.DB {
		return db
	})

	c.Provide(logger.NewLogger)
	c.Provide(token.NewJwtService)
	c.Provide(repository.NewUserRepository)
	c.Provide(service.NewUserService)
	c.Provide(middlewares.NewAuthMiddleware)
	c.Provide(handler.NewUserHandler)
}
