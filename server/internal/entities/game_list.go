package entities

import (
	"time"

	"github.com/google/uuid"
)

const (
	DefaultListName = "Watchlist"
)

type GameList struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	IsPublic  bool      `gorm:"type:bool;default:true"`
	IsDefault bool      `gorm:"type:bool;not null"`
	CreatedAt time.Time `gorm:"autoCreatedTime"`
	UpdatedAt time.Time `gorm:"autoUpdatedTime"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Game   []Game    `gorm:"many2many:list_items"`
}
