package rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	null "gopkg.in/guregu/null.v4"
)

// SearchActivePackageOutput output
type SearchActivePackageOutput struct {
	ID                   uuid.UUID                  `json:"id"`
	Questionnaire        model.Questionnaire        `json:"questionnaire"`
	IndicationCategories model.IndicationCategories `json:"indication_categories"`
	Name                 string                     `json:"name"`
}

// QuestionnaireGrade components of what is considered grade or score from each submitted questionnaire
type QuestionnaireGrade struct {
	Detail     model.ResultDetail       `json:"detail"`
	Indication model.IndicationCategory `json:"indication"`
	Total      int                      `json:"total"`
}

// SubmitQuestionnaireOutput output
type SubmitQuestionnaireOutput struct {
	ResultID  uuid.UUID          `json:"result_id"`
	Grade     QuestionnaireGrade `json:"grade"`
	ChildID   uuid.UUID          `json:"child_id,omitempty"`
	CreatedBy uuid.UUID          `json:"created_by,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
}

// GetChildStatOutput output
type GetChildStatOutput struct {
}

// SearchQUestionnaireResultsOutput output
type SearchQUestionnaireResultsOutput struct {
	ID        uuid.UUID          `json:"id"`
	PackageID uuid.UUID          `json:"package_id"`
	ChildID   uuid.UUID          `json:"child_id"`
	CreatedBy uuid.UUID          `json:"created_by"`
	Answer    model.AnswerDetail `json:"answer"`
	Result    model.ResultDetail `json:"result"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt null.Time          `json:"deleted_at,omitempty"`
}

// GetMyQUestionnaireResultsOutput output
type GetMyQUestionnaireResultsOutput struct {
	TODO []any `json:"todo"`
}

// SignupOutput output
type SignupOutput struct {
	Message string `json:"message" example:"confirmation link sent to your email"`
}

// VerifyAccountOutput output
type VerifyAccountOutput struct {
	Message string `json:"message" example:"your account is now activated and can be used"`
}

// InitResetPasswordOutput output
type InitResetPasswordOutput struct {
	Message string `json:"message"`
}

// ResetPasswordOutput output
type ResetPasswordOutput struct {
	Message string `json:"message"`
}

// LoginOutput output
type LoginOutput struct {
	Token string `json:"token"`
}

// RegisterChildOutput output
type RegisterChildOutput struct {
	ID uuid.UUID `json:"id"`
}

// UpdateChildernOutput output
type UpdateChildernOutput struct {
	Message string `json:"message"`
}

// GetMyChildernOutput output
type GetMyChildernOutput struct {
	ID             uuid.UUID   `json:"id"`
	ParentUserID   uuid.UUID   `json:"parent_user_id"`
	ParentUserName string      `json:"parent_user_name"`
	DateOfBirth    time.Time   `json:"date_of_birth"`
	Gender         bool        `json:"gender"`
	Name           string      `json:"name"`
	GuardianName   null.String `json:"guardian_name" swaggertype:"string"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// CreatePackageOutput output
type CreatePackageOutput struct {
	ID uuid.UUID `json:"id"`
}

// UpdatePackageOutput output
type UpdatePackageOutput struct {
	Message string
}

// ActivationPackageOutput output
type ActivationPackageOutput struct {
	Message string
}

// CreateATECTemplateOutput output
type CreateATECTemplateOutput struct {
	ID uuid.UUID `json:"id"`
}

// UpdateATECTemplateOutput output
type UpdateATECTemplateOutput struct {
	ID uuid.UUID `json:"id"`
}

// GetATECQuestionnaireOutput output
type GetATECQuestionnaireOutput struct {
	ID            uuid.UUID           `json:"id"`
	Questionnaire model.Questionnaire `json:"questionnaire"`
	Name          string              `json:"name"`
}

// SearchChildernOutput output
type SearchChildernOutput struct {
	ID           uuid.UUID   `json:"id"`
	ParentUserID uuid.UUID   `json:"parent_user_id"`
	DateOfBirth  time.Time   `json:"date_of_birth"`
	Gender       bool        `json:"gender"`
	Name         string      `json:"name"`
	GuardianName null.String `json:"guardian_name" swaggertype:"string"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// ResendVerificationOutput output
type ResendVerificationOutput struct {
	Message string `json:"message"`
}
