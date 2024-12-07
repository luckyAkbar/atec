// Package model contains datatyped related to database model and / or core system data type
package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Child represent childern table on database
type Child struct {
	ID           uuid.UUID
	ParentUserID uuid.UUID
	DateOfBirth  time.Time
	Gender       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}
