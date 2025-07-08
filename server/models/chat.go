package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	User1ID   uuid.UUID `json:"user1_id" gorm:"not null"`
	User2ID   uuid.UUID `json:"user2_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
