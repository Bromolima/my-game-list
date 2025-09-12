package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/Bromolima/my-game-list/config"
	"github.com/Bromolima/my-game-list/database"
	"github.com/Bromolima/my-game-list/internal/http/routes"
	"github.com/Bromolima/my-game-list/internal/injector"
	validation "github.com/Bromolima/my-game-list/internal/validation"
	_ "github.com/Bromolima/my-game-list/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func main() {
	slog.Info("starting server", slog.String("port", config.Env.ApiPort))
	if err := config.LoadEnvironment(); err != nil {
		log.Fatal(err)
	}

	c := dig.New()
	e := echo.New()
	v := validation.NewCustomValidator()

	validation.SetupTranslations(v)
	e.Validator = v

	db, err := database.SetupPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	injector.SetupDependecies(c, db)

	if err := routes.SetupRoutes(e, c); err != nil {
		log.Fatal(err)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Env.ApiPort)))
}
