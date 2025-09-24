package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Password  string    `gorm:"type:varchar(100);not null"`
	Username  string    `gorm:"type:varchar(100);not null"`
	AvatarURL string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	RoleID    uint
}
