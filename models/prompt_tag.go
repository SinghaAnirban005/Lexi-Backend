package models

import (
	"github.com/google/uuid"
)

type PromptTag struct {
	PromptID uuid.UUID `gorm:"type:uuid;primaryKey" json:"prompt_id"`
	TagID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"tag_id"`

	Prompt Prompt `gorm:"foreignKey:PromptID;constraint:OnDelete:CASCADE;"`
	Tag    Tag    `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE;"`
}
