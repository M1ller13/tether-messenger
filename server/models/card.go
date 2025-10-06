package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Card struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Position    int            `json:"position" gorm:"not null"`
	Color       string         `json:"color" gorm:"default:'#FFFFFF'"`
	ColumnID    uuid.UUID      `json:"column_id" gorm:"type:uuid;not null"`
	AssigneeID  *uuid.UUID     `json:"assignee_id" gorm:"type:uuid"`
	CreatedByID uuid.UUID      `json:"created_by_id" gorm:"type:uuid;not null"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// CRM Fields
	LeadName     string  `json:"lead_name" gorm:"type:varchar(255)"`
	ContactEmail string  `json:"contact_email" gorm:"type:varchar(255)"`
	ContactPhone string  `json:"contact_phone" gorm:"type:varchar(50)"`
	Company      string  `json:"company" gorm:"type:varchar(255)"`
	Value        float64 `json:"value" gorm:"type:decimal(15,2)"`
	Priority     string  `json:"priority" gorm:"default:'medium'"` // 'low', 'medium', 'high', 'urgent'
	Status       string  `json:"status" gorm:"default:'new'"`      // 'new', 'contacted', 'qualified', 'proposal', 'negotiation', 'closed-won', 'closed-lost'

	// Relations
	Column    Column `json:"column,omitempty" gorm:"foreignKey:ColumnID"`
	Assignee  *User  `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	CreatedBy User   `json:"created_by" gorm:"foreignKey:CreatedByID"`
}
