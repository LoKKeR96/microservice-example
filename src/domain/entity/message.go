package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Text      string    `gorm:"size:255"`
	UUID      uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Message) TableName() string {
	return "messages"
}
