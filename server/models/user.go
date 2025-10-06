package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email         string         `json:"email" gorm:"unique;not null"`
	Password      string         `json:"-" gorm:"not null"` // Hidden from JSON
	Username      string         `json:"username" gorm:"unique;not null"`
	DisplayName   string         `json:"display_name" gorm:"not null"`
	Bio           string         `json:"bio" gorm:"type:varchar(300)"`
	AvatarURL     string         `json:"avatar_url"`
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	LastSeen      *time.Time     `json:"last_seen"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
