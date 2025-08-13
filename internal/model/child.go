// Package model contains datatyped related to database model and / or core system data type
package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Child represent childern table on database
type Child struct {
	ID           uuid.UUID `gorm:"default:uuid_generate_v4()"`
	ParentUserID uuid.UUID
	DateOfBirth  time.Time
	Gender       bool
	Name         string
	GuardianName sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}
