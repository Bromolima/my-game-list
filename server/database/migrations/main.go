package main

import (
	"log"
	"log/slog"

	"github.com/Bromolima/my-game-list/config"
	"github.com/Bromolima/my-game-list/database"
	"github.com/Bromolima/my-game-list/internal/entities"
	_ "github.com/Bromolima/my-game-list/logger"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnvironment()

	db, err := database.SetupPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&entities.Access{},
		&entities.Role{},
		&entities.User{},
		&entities.Game{},
		&entities.GameList{},
		&entities.ListItem{},
	); err != nil {
		log.Fatal(err)
	}

	if err := setupAccessess(db); err != nil {
		log.Fatal(err)
	}

	slog.Info("database migration completed successfully")
}

func setupAccessess(db *gorm.DB) error {
	db.Create(&entities.Access{AccessType: entities.ReadAccess})
	db.Create(&entities.Access{AccessType: entities.UpdateAccess})
	db.Create(&entities.Access{AccessType: entities.CreateAcess})
	db.Create(&entities.Access{AccessType: entities.DeleteAcess})

	var read, update, create, delete entities.Access
	db.First(&read, "access_type = ?", entities.ReadAccess)
	db.First(&update, "access_type = ?", entities.UpdateAccess)
	db.First(&create, "access_type = ?", entities.CreateAcess)
	db.First(&delete, "access_type = ?", entities.DeleteAcess)

	adminRole := entities.Role{
		ID:     entities.RoleAdminID,
		Name:   entities.RoleAdminName,
		Access: []entities.Access{read, update, create, delete},
	}

	userRole := entities.Role{
		ID:     entities.RoleUserID,
		Name:   entities.RoleUserName,
		Access: []entities.Access{read},
	}

	if err := db.Create(&userRole).Error; err != nil {
		return err
	}

	if err := db.Create(&adminRole).Error; err != nil {
		return err
	}

	return nil
}
