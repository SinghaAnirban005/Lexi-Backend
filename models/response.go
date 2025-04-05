package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Response struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Response  string    `json:"response"`
	PromptID  uuid.UUID `gorm:"type:uuid" json:"prompt_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Prompt    Prompt    `gorm:"foreignKey:PromptID;constraint:OnDelete:CASCADE;" json:"prompt"`
}

func (r *Response) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
