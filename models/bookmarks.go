package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bookmark struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ConversationID uuid.UUID `gorm:"type:uuid" json:"conversation_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`

	User         User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Conversation Conversation `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"conversation"`
}

func (b *Bookmark) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
