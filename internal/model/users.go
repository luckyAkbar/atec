package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Roles string

const (
	RolesAdmin Roles = "admin"
	RoleUser   Roles = "user"
)

// User represent users table on database
type User struct {
	ID        uuid.UUID `gorm:"default:uuid_generate_v4()"`
	Email     string
	Password  string `json:"-"`
	Username  string
	IsActive  bool
	Roles     Roles
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
