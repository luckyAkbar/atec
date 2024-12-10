package model

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Package represent packages table on database
type Package struct {
	ID            uuid.UUID `gorm:"default:uuid_generate_v4()"`
	CreatedBy     uuid.UUID
	Questionnaire Questionnaire
	Name          string
	IsActive      bool
	IsLocked      bool
	CreatedAt     time.Time `gorm:"default:now()"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

// AnswerOption each singular answer option for a given question along with its detail
type AnswerOption struct {
	ID          int    `json:"id" validate:"required"`
	Description string `json:"description" validate:"required"`
	Score       int    `json:"score" validate:"min=0"`
}

// ChecklistGroup is each individual question and answer in a given group
type ChecklistGroup struct {
	// CustomName can be used to be displayed to user. if empty, the one from
	// the template will be used
	CustomName string         `json:"custom_name" validate:"required"`
	Questions  []string       `json:"questions" validate:"required,min=1"`
	Options    []AnswerOption `json:"options" validate:"required,min=1,unique=ID,unique=Score"`
}

// Validate validate ChecklistGroup based on its rules on the struct's tags
// and also ensuring all options have complete range of Score
func (cg ChecklistGroup) Validate() error {
	if err := common.Validator.Struct(cg); err != nil {
		return err
	}

	// to ensure each options score cover from 0 to len(cg.Options)
	for i := range len(cg.Options) {
		found := false

		for _, opt := range cg.Options {
			if opt.Score == i {
				found = true

				break
			}
		}

		if !found {
			return fmt.Errorf(
				"missing option score %d for checklist group name %s",
				i,
				cg.CustomName,
			)
		}
	}

	return nil
}

// Questionnaire represent all the ATEC questionnaire structure
type Questionnaire map[int]ChecklistGroup

// Validate ensure the supplied Questionnaire match with the ATEC template
func (q Questionnaire) Validate() error {
	for j, template := range DefaultATECTemplate.SubTest {
		cg, ok := q[j]
		if !ok {
			return fmt.Errorf("questionnaire group number %d for %s not found", j+1, template.Name)
		}

		if err := cg.Validate(); err != nil {
			return err
		}

		if len(cg.Options) != template.OptionCount {
			return fmt.Errorf(
				"questionnaire group number %d for %s expecting %d number of options, but got %d",
				j+1,
				template.Name,
				template.OptionCount,
				len(cg.Options),
			)
		}

		if len(cg.Questions) != template.QuestionCount {
			return fmt.Errorf(
				"questionnaire group number %d for %s expecting %d number of questions, but got %d",
				j+1,
				template.Name,
				template.QuestionCount,
				len(cg.Questions),
			)
		}
	}

	return nil
}

// Value implements Valuer/Scanner interface
func (q Questionnaire) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}

// Scan implements Valuer/Scanner interface
func (q *Questionnaire) Scan(_ context.Context, _ *schema.Field, _ reflect.Value, dbValue interface{}) error {
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

	if err := json.Unmarshal(bytes, q); err != nil {
		return err
	}

	return nil
}
