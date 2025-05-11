package rest

import (
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
)

// SubmitQuestionnaireInput input
type SubmitQuestionnaireInput struct {
	PackageID uuid.UUID          `json:"package_id"`
	ChildID   uuid.UUID          `json:"child_id"`
	Answers   model.AnswerDetail `validate:"required" json:"answers"`
}

// SearchQUestionnaireResultsInput input
type SearchQUestionnaireResultsInput struct {
	ResultID    uuid.UUID `json:"result_id" query:"result_id"`
	PackageID   uuid.UUID `json:"package_id" query:"package_id"`
	ChildID     uuid.UUID `json:"child_id" query:"child_id"`
	CreatedByID uuid.UUID `json:"created_by_id" query:"created_by_id"`
	Limit       int       `json:"limit" query:"limit" validate:"min=1,max=100"`
	Offset      int       `json:"offset" query:"offset" validate:"min=0"`
}

// GetMyQUestionnaireResultsInput input
type GetMyQUestionnaireResultsInput struct {
	Limit  int `json:"limit" query:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" query:"offset" validate:"min=0"`
}

// SignupInput input
type SignupInput struct {
	Email    string `json:"email" validate:"required,email" example:"string@string.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
	Username string `json:"username" validate:"required" example:"username"`
}

// VerifyAccountInput input
type VerifyAccountInput struct {
	ValidationToken string `query:"validation_token" validate:"required"`
}

// InitResetPasswordInput input
type InitResetPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordInput input
type ResetPasswordInput struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// LoginInput input
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterChildInput input
type RegisterChildInput struct {
	Name        string `json:"name" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required" example:"2001-11-29 (YYYY-MM-DD)"`
	Gender      bool   `json:"gender" example:"true"`
}

// UpdateChildernInput input
type UpdateChildernInput struct {
	ChildID     uuid.UUID `json:"-" param:"child_id"`
	Name        *string   `json:"name" validate:"required"`
	DateOfBirth string    `json:"date_of_birth" validate:"required" example:"2001-11-29 (YYYY-MM-DD)"`
	Gender      *bool     `json:"gender" example:"true"`
}

// CreatePackageInput input
type CreatePackageInput struct {
	PackageName             string                        `json:"package_name" validate:"required"`
	Quesionnaire            model.Questionnaire           `json:"questionnaire" validate:"required"`
	IndicationCategories    model.IndicationCategories    `json:"indication_categories" validate:"required"`
	ImageResultAttributeKey model.ImageResultAttributeKey `json:"image_result_attribute_key" validate:"required"`
}

// UpdatePackageInput input
type UpdatePackageInput struct {
	PackageID uuid.UUID `json:"-" param:"package_id"`

	CreatePackageInput
}

// ActivationPackageInput input
type ActivationPackageInput struct {
	Status    bool      `json:"status"`
	PackageID uuid.UUID `json:"-" param:"package_id"`
}

// SDTemplateSubGroupDetail input
type SDTemplateSubGroupDetail struct {
	Name              string `json:"name" validate:"required"`
	QuestionCount     int    `json:"question_count" validate:"required,min=1"`
	AnswerOptionCount int    `json:"answer_option_count" validate:"required,min=2"`
}

// CreateATECTemplateInput input
type CreateATECTemplateInput struct {
	Name                   string                     `json:"name" validate:"required,max=255"`
	IndicationThreshold    int                        `json:"indication_threshold" validate:"required,min=0"`
	PositiveIndiationText  string                     `json:"positive_indication_text" validate:"required"`
	NegativeIndicationText string                     `json:"negative_indication_text" validate:"required"`
	SubGroupDetails        []SDTemplateSubGroupDetail `json:"sub_group_details" validate:"min=1,dive"`
}

// GetATECQuestionnaireInput input
type GetATECQuestionnaireInput struct {
	PackageID uuid.UUID `query:"package_id"`
}

// DeletePackageInput input
type DeletePackageInput struct {
	PackageID uuid.UUID `param:"package_id"`
}

// GetMyChildrenInput input
type GetMyChildrenInput struct {
	Limit  int `query:"limit" validate:"min=1"`
	Offset int `query:"offset" validate:"min=0"`
}

// SearchChildrenInput input
type SearchChildrenInput struct {
	ParentUserID *uuid.UUID `json:"parent_user_id" query:"parent_user_id"`
	Name         *string    `query:"name"`
	Gender       *bool      `query:"gender"`
	Limit        int        `query:"limit" validate:"min=1" example:"1"`
	Offset       int        `query:"offset" validate:"min=0"`
}

// DownloadQuestionnaireResultInput input
type DownloadQuestionnaireResultInput struct {
	ResultID uuid.UUID `param:"result_id" validate:"required"`
}

// GetChildStatInput input
type GetChildStatInput struct {
	ChildID uuid.UUID `param:"child_id" validate:"required"`
}

// ResendVerificationInput input
type ResendVerificationInput struct {
	Email string `json:"email" validate:"required,email"`
}

// DeleteAccountInput input
type DeleteAccountInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
