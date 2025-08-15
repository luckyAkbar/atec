package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Roles represent database's enum for roles
type Roles string

// known roles
const (
	RolesAdministrator Roles = "administrator"
	RolesParent        Roles = "parent"
	RolesTherapist     Roles = "therapist"
)

// User represent users table on database
type User struct {
	ID          uuid.UUID `gorm:"default:uuid_generate_v4()"`
	Email       string
	Password    string `json:"-"`
	Username    string
	IsActive    bool
	Roles       Roles
	PhoneNumber sql.NullString
	Address     sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
