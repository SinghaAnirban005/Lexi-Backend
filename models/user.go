package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	FullName      string         `json:"full_name"`
	Username      string         `gorm:"uniqueIndex" json:"username"`
	Email         string         `gorm:"uniqueIndex" json:"email"`
	Password      string         `json:"-"`
	Conversations []Conversation `gorm:"foreignKey:OwnerID" json:"conversations"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
