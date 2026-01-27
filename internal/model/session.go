package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"session_id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null" json:"account_id"`
	Account   Account   `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"-"`

	Data      json.RawMessage `gorm:"type:jsonb" json:"data,omitempty"`
	IPAddress string          `gorm:"type:inet" json:"ip_address"`
	UserAgent string          `json:"user_agent"`

	CreatedAt  time.Time  `gorm:"type:timestamptz;default:now()" json:"created_at"`
	LastSeenAt *time.Time `gorm:"type:timestamptz" json:"last_seen_at,omitempty"`
	RevokedAt  *time.Time `gorm:"type:timestamptz" json:"revoked_at,omitempty"`
}

func (s *Session) TableName() string { return "sessions" }
