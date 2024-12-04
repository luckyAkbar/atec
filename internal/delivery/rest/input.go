package rest

import (
	"time"

	"github.com/google/uuid"
)

type SearchChildernOutput struct {
	Childern []any `json:"childern"`
}

type SubmitQuestionnaireInput struct {
	PackageID uuid.UUID `json:"package_id"`
	ChildID   uuid.UUID `json:"child_id"`
	TODO      []any     `json:"todo"`
}

type SearchQUestionnaireResultsInput struct {
	ResultID uuid.UUID `query:"result_id"`
	UserID   uuid.UUID `query:"user_id"`
	ChildID  uuid.UUID `query:"child_id"`
	Limit    int       `query:"limit"`
	Offset   int       `query:"offset"`
}

type GetMyQUestionnaireResultsInput struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type SignupInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"username" validate:"required"`
}

type VerifyAccountInput struct {
	ValidationToken string
}

type InitResetPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterChildInput struct {
	Name        string    `json:"name" validate:"required"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required" example:"2001-11-29 (YYYY-MM-DD)"`
}

type UpdateChildernInput struct {
	RegisterChildInput
}

type AnswerAndValue struct {
	Text  string `json:"text" validate:"required"`
	Value int    `json:"value" validate:"required,min=1"`
}

type QuestionAndAnswers struct {
	Question        string           `json:"question" validate:"required"`
	AnswersAndValue []AnswerAndValue `json:"answer_and_value" validate:"required,min=1,unique=Value,dive"`
}

type SubGroupDetail struct {
	Name                   string               `json:"name" validate:"required"`
	QuestionAndAnswerLists []QuestionAndAnswers `json:"question_and_answer_lists" validate:"required,min=1,dive"`
}

type CreatePackageInput struct {
	PackageName     string           `json:"package_name" validate:"required"`
	SubGroupDetails []SubGroupDetail `json:"sub_group_details" validate:"required,min=1,unique=Name,dive"`
}

type UpdatePackageInput struct {
	CreatePackageInput
}

type ActivationPackageInput struct {
	Status bool `json:"status"`
}

type SDTemplateSubGroupDetail struct {
	Name              string `json:"name" validate:"required"`
	QuestionCount     int    `json:"question_count" validate:"required,min=1"`
	AnswerOptionCount int    `json:"answer_option_count" validate:"required,min=2"`
}

type CreateATECTemplateInput struct {
	Name                   string                     `json:"name" validate:"required,max=255"`
	IndicationThreshold    int                        `json:"indication_threshold" validate:"required,min=0"`
	PositiveIndiationText  string                     `json:"positive_indication_text" validate:"required"`
	NegativeIndicationText string                     `json:"negative_indication_text" validate:"required"`
	SubGroupDetails        []SDTemplateSubGroupDetail `json:"sub_group_details" validate:"min=1,dive"`
}

type UpdateATECTemplateInput struct {
	CreateATECTemplateInput
}

type ActivateTemplateInput struct {
	Status bool `json:"status"`
}

type GetATECQuestionnaireInput struct {
	PackageID uuid.UUID `query:"package_id,omitempty"`
}
