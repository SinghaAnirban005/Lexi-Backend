package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Prompt struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;" json:"id"`
	PromptTitle    string       `json:"prompt_title"`
	ConversationID uuid.UUID    `gorm:"type:uuid" json:"conversation_id"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE;" json:"conversation"`
	Responses      []Response   `gorm:"foreignKey:PromptID" json:"responses"`
}

func (p *Prompt) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
