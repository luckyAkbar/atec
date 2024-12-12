package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Result represent results table on database
type Result struct {
	ID        uuid.UUID
	PackageID uuid.UUID
	CreatedBy uuid.UUID `gorm:"default:null"`
	Answer    AnswerDetail
	Result    ResultDetail
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

// AnswerDetail represent each checklisted option from the questionnaire.
// The first key (int) is the subtest id. each subtest will have
// map with key in int and the value also in int. That map
// represent each question (key) and checklisted option (value)
type AnswerDetail map[int]map[int]int

// ResultDetail will contain each result from the questionnaire's group
type ResultDetail map[int]int
