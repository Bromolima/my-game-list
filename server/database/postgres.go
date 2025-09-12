package database

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Bromolima/my-game-list/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupPostgresConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Env.DB.Host,
		config.Env.DB.User,
		config.Env.DB.Password,
		config.Env.DB.Name,
		config.Env.DB.Port,
	)

	l := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: l,
	})

	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("database connection successful")
	return db, nil
}
