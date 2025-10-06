package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Workspace struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Type        string         `json:"type" gorm:"not null;default:'team'"` // 'team', 'organization'
	OwnerID     uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	Slug        string         `json:"slug" gorm:"unique;not null"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Owner   User              `json:"owner" gorm:"foreignKey:OwnerID"`
	Boards  []Board           `json:"boards,omitempty" gorm:"foreignKey:WorkspaceID"`
	Members []WorkspaceMember `json:"members,omitempty" gorm:"foreignKey:WorkspaceID"`
}

type WorkspaceMember struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkspaceID uuid.UUID      `json:"workspace_id" gorm:"type:uuid;not null"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Role        string         `json:"role" gorm:"not null;default:'member'"` // 'owner', 'admin', 'member'
	JoinedAt    time.Time      `json:"joined_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Workspace Workspace `json:"workspace,omitempty" gorm:"foreignKey:WorkspaceID"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
}
