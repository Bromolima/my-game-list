package injector

import (
	"github.com/Bromolima/my-game-list/internal/entities"
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

	c.Provide(repository.NewRoleRepository)
	c.Provide(repository.NewPageRepository[entities.Game])
	c.Provide(repository.NewPageRepository[entities.User])
	c.Provide(repository.NewUserRepository)
	c.Provide(repository.NewGameRepository)
	c.Provide(repository.NewListItemRepository)

	c.Provide(token.NewJwtService)
	c.Provide(service.NewGameListService)
	c.Provide(service.NewListItemService)
	c.Provide(service.NewGameService)
	c.Provide(service.NewUserService)

	c.Provide(middlewares.NewAuthMiddleware)

	c.Provide(handler.NewGameListHandler)
	c.Provide(handler.NewGameHandler)
	c.Provide(handler.NewUserHandler)
}
