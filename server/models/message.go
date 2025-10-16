package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChatID   uuid.UUID `json:"chat_id" gorm:"not null"`
	SenderID uuid.UUID `json:"sender_id" gorm:"not null"`
	// Plaintext content for non-E2EE chats (legacy). Should be empty when E2EE is used.
	Content string `json:"content"`
	// E2EE payload fields
	Ciphertext   string    `json:"ciphertext" gorm:"type:text"`
	Nonce        string    `json:"nonce" gorm:"type:varchar(64)"`
	Alg          string    `json:"alg" gorm:"type:varchar(64)"`
	EphemeralPub string    `json:"ephemeral_pub" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
	IsRead       bool      `json:"is_read"`
}
