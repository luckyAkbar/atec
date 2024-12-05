package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Childern represent childern table on database
type Child struct {
	ID           uuid.UUID
	ParentUserID uuid.UUID
	DateOfBirth  time.Time
	Gender       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}
