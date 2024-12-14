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
	ID                   uuid.UUID `gorm:"default:uuid_generate_v4()"`
	CreatedBy            uuid.UUID
	Questionnaire        Questionnaire
	IndicationCategories IndicationCategories
	Name                 string
	IsActive             bool
	IsLocked             bool
	CreatedAt            time.Time `gorm:"default:now()"`
	UpdatedAt            time.Time `gorm:"default:now()"`
	DeletedAt            gorm.DeletedAt
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

// IndicationCategories is the template to decide what indication for a given score.
type IndicationCategories []IndicationCategory

// GetIndicationCategoryByScore will return the indication category based on score.
func (ic IndicationCategories) GetIndicationCategoryByScore(score int) IndicationCategory {
	for _, category := range ic {
		if category.IsInTheRange(score) {
			return category
		}
	}

	// this is a fallback mechanism. if the score is not found, return an invalid category.
	// eventhough this should never happen, it is better to have a fallback mechanism instead of returning
	// dangerous nil value
	return IndicationCategory{
		MinimumScore: 0,
		MaximumScore: 0,
		Name:         "invalid value",
		Detail:       "invalid indication",
	}
}

// Validate will validate given IndicationCategories to cover entire possible range
// of ATEC score. Also, to ensure that no overlapping / missing gap on each categories.
func (ic IndicationCategories) Validate() error {
	minScore := DefaultATECTemplate.MinimumPossibleScore
	maxScore := DefaultATECTemplate.MaximumPossibleScore

	for score := minScore; score <= maxScore; score++ {
		matchCount := 0

		for _, category := range ic {
			if category.IsInTheRange(score) {
				matchCount++
			}

			if err := category.Validate(); err != nil {
				return err
			}
		}

		if matchCount != 1 {
			return fmt.Errorf("score: %d has no matching categories or overlapping on the defined categories", score)
		}
	}

	return nil
}

// Value implements Valuer/Scanner interface
func (ic IndicationCategories) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}

// Scan implements Valuer/Scanner interface
func (ic *IndicationCategories) Scan(_ context.Context, _ *schema.Field, _ reflect.Value, dbValue interface{}) error {
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

	if err := json.Unmarshal(bytes, ic); err != nil {
		return err
	}

	return nil
}

// IndicationCategory is the detailed structure composing a category to decide indication based on a given score
type IndicationCategory struct {
	MinimumScore int    `json:"minimum_score" validate:"min=0"`
	MaximumScore int    `json:"maximum_score" validate:"required,min=1"`
	Name         string `json:"name" validate:"required"`
	Detail       string `json:"detail" validate:"required"`
}

// IsInTheRange will check if the given score is in the range of this category
func (ic IndicationCategory) IsInTheRange(score int) bool {
	return score >= ic.MinimumScore && score <= ic.MaximumScore
}

// Validate will validate the given IndicationCategory
func (ic IndicationCategory) Validate() error {
	if err := common.Validator.Struct(ic); err != nil {
		return err
	}

	if ic.MinimumScore < DefaultATECTemplate.MinimumPossibleScore {
		return fmt.Errorf("minimum score %d is less than minimum possible score %d", ic.MinimumScore, DefaultATECTemplate.MinimumPossibleScore)
	}

	if ic.MaximumScore > DefaultATECTemplate.MaximumPossibleScore {
		return fmt.Errorf("maximum score %d is greater than maximum possible score %d", ic.MaximumScore, DefaultATECTemplate.MaximumPossibleScore)
	}

	return nil
}
