package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Package represent packages table on database
type Package struct {
	ID        uuid.UUID
	CreatedBy uuid.UUID

	// TODO will be defined later
	Questionnaire Questionnaire
	Name          string
	IsActive      bool
	IsLocked      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
}

// TODO use this
type AnswerOption struct {
	ID          int
	Description string
	Score       int
}

// TODO use this
type QuestionAndOptions struct {
	QUestion string
	Options  []AnswerOption
}

type Questionnaire struct {
	ChecklistGroups []ChecklistGroup
}

type Checklists map[int]QuestionAndOptions

type ChecklistGroup struct {
	SubtestID int

	// CustomName can be used to be displayed to user. if empty, the one from
	// the template will be used
	CustomName string
	Checklists Checklists
}
