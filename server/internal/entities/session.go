package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Token     string    `gorm:"type:text"`
	IPAdress  string    `gorm:"type:text"`
	UserAgent string    `gorm:"type:text"`
	ExpiresAt time.Time `gorm:"type:timestamp"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:Cascade"`
}
