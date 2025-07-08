package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Phone       string         `json:"phone" gorm:"unique;not null"`
	Username    string         `json:"username" gorm:"unique;not null"`
	DisplayName string         `json:"display_name" gorm:"not null"`
	Bio         string         `json:"bio" gorm:"type:varchar(300)"`
	AvatarURL   string         `json:"avatar_url"`
	LastSeen    *time.Time     `json:"last_seen"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
