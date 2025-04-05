package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID      uuid.UUID   `gorm:"type:uuid;primaryKey;" json:"id"`
	TagName string      `gorm:"uniqueIndex" json:"tag_name"`
	Prompts []PromptTag `gorm:"foreignKey:TagID" json:"prompts"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}
