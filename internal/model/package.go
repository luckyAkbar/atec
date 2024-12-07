package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Package represent packages table on database
type Package struct {
	ID            uuid.UUID
	CreatedBy     uuid.UUID
	Questionnaire Questionnaire
	Name          string
	IsActive      bool
	IsLocked      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
}

// AnswerOption each singular answer option for a given question along with its detail
type AnswerOption struct {
	ID          int
	Description string
	Score       int
}

// QuestionAndOptions structure for a question and each answer options
type QuestionAndOptions struct {
	QUestion string
	Options  []AnswerOption
}

// Questionnaire represent all the ATEC questionnaire structure
type Questionnaire struct {
	ChecklistGroups []ChecklistGroup
}

// Checklists is representing a group of questions and answers for a given sub test / group
type Checklists map[int]QuestionAndOptions

// ChecklistGroup is each individual question and answer in a given group
type ChecklistGroup struct {
	SubtestID int

	// CustomName can be used to be displayed to user. if empty, the one from
	// the template will be used
	CustomName string
	Checklists Checklists
}
