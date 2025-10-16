package models

import (
	"time"

	"github.com/google/uuid"
)

// DeviceKey stores public identity and signed prekey info for a user's device.
// Only PUBLIC materials are stored here; private keys never touch the server.
type DeviceKey struct {
	ID                    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID                uuid.UUID `json:"user_id" gorm:"type:uuid;index;not null"`
	DeviceID              string    `json:"device_id" gorm:"type:varchar(128);not null;index:idx_user_device,unique"`
	IdentityKeyPublic     string    `json:"identity_key_public" gorm:"type:text;not null"`
	SignedPreKeyPublic    string    `json:"signed_prekey_public" gorm:"type:text;not null"`
	SignedPreKeySignature string    `json:"signed_prekey_signature" gorm:"type:text;not null"`
	Active                bool      `json:"active" gorm:"default:true"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// OneTimePreKey stores a pool of one-time prekeys for a device.
// Each key may be served at most once to a peer to start a session.
type OneTimePreKey struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DeviceKeyID uuid.UUID  `json:"device_key_id" gorm:"type:uuid;index;not null"`
	KeyID       int        `json:"key_id" gorm:"not null;index:idx_device_keyid,unique"`
	PublicKey   string     `json:"public_key" gorm:"type:text;not null"`
	Used        bool       `json:"used" gorm:"default:false;index"`
	CreatedAt   time.Time  `json:"created_at"`
	UsedAt      *time.Time `json:"used_at"`
}
