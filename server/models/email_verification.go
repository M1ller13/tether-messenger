package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"not null"`
	Token     string    `json:"token" gorm:"unique;not null"`
	Type      string    `json:"type" gorm:"not null"` // "signup" or "password_reset"
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	Used      bool      `json:"used" gorm:"default:false"`
}
