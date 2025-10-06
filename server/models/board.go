package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Type        string         `json:"type" gorm:"not null;default:'personal'"` // 'personal', 'team', 'crm'
	OwnerID     uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	WorkspaceID *uuid.UUID     `json:"workspace_id" gorm:"type:uuid"` // null for personal boards
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	Color       string         `json:"color" gorm:"default:'#3B82F6'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Owner     User      `json:"owner" gorm:"foreignKey:OwnerID"`
	Workspace Workspace `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID"`
	Columns   []Column  `json:"columns,omitempty" gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
}
