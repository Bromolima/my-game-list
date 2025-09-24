package entities

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Genre       string    `gorm:"type:varchar(100);not null"`
	Developer   string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`
	Rating      float32   `gorm:"type:real;not null;default:0"`
	ImageURL    string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Game        []Game    `gorm:"many2many:list_items"`
}
