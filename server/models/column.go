package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Column struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null"`
	Position  int            `json:"position" gorm:"not null"`
	Color     string         `json:"color" gorm:"default:'#6B7280'"`
	BoardID   uuid.UUID      `json:"board_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Board Board  `json:"board,omitempty" gorm:"foreignKey:BoardID"`
	Cards []Card `json:"cards,omitempty" gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE"`
}
