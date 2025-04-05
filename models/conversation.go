package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	OwnerID   uuid.UUID `gorm:"type:uuid" json:"owner_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	User      User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Prompts   []Prompt  `gorm:"foreignKey:ConversationID" json:"prompts"`
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
