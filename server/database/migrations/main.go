package main

import (
	"log"
	"log/slog"

	"github.com/Bromolima/my-game-list/config"
	"github.com/Bromolima/my-game-list/database"
	"github.com/Bromolima/my-game-list/internal/models"
	_ "github.com/Bromolima/my-game-list/logger"
)

func main() {
	config.LoadEnvironment()

	db, err := database.SetupPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	slog.Info("database migration completed successfully")
}
