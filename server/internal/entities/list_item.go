package entities

import (
	"time"

	"github.com/google/uuid"
)

const (
	ListItemStatusWant         = "WANT_TO_PLAY"
	ListItemStatusPlaying      = "PLAYING"
	ListItemStatusDoneFinished = "FINISHED"
)

type ListItem struct {
	GameListID uuid.UUID `gorm:"primaryKey;type:uuid"`
	GameID     uuid.UUID `gorm:"primaryKey;type:uuid"`
	Status     string    `gorm:"type:varchar(100)"`
	Rating     float32   `gorm:"type:numeric"`
	CreatedAt  time.Time `gorm:"autoCreatedTime"`
	UpdatedAt  time.Time `gorm:"autoUpdatedTime"`
}
