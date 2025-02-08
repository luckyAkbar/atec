package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/schema"
)

// Result represent results table on database
type Result struct {
	ID        uuid.UUID `gorm:"default:uuid_generate_v4()"`
	PackageID uuid.UUID
	ChildID   uuid.UUID `gorm:"default:null"`
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

// Value implements Valuer / Scanner interface to be compatible as JSONB field on postgres
func (ad AnswerDetail) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}

// Scan implements Valuer / Scanner interface to be compatible as JSONB field on postgres
func (ad *AnswerDetail) Scan(_ context.Context, _ *schema.Field, _ reflect.Value, dbValue interface{}) error {
	if dbValue == nil {
		return nil
	}

	var bytes []byte
	switch v := dbValue.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
	}

	if err := json.Unmarshal(bytes, ad); err != nil {
		return err
	}

	return nil
}

// SubtestGrade represent each subtest's grade from a particular subtest group
type SubtestGrade struct {
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

// ResultDetail will contain each result from the questionnaire's group
type ResultDetail map[int]SubtestGrade

// Value implements Valuer / Scanner interface to be compatible as JSONB field on postgres
func (rd ResultDetail) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}

// Scan implements Valuer / Scanner interface to be compatible as JSONB field on postgres
func (rd *ResultDetail) Scan(_ context.Context, _ *schema.Field, _ reflect.Value, dbValue interface{}) error {
	if dbValue == nil {
		return nil
	}

	var bytes []byte
	switch v := dbValue.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
	}

	if err := json.Unmarshal(bytes, rd); err != nil {
		return err
	}

	return nil
}

// CountTotalScore will count the total score of each subtest
func (rd ResultDetail) CountTotalScore() int {
	total := 0
	for _, v := range rd {
		total += v.Grade
	}

	return total
}
