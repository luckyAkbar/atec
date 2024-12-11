package rest

import "github.com/google/uuid"

// SearchActivePackageOutput output
type SearchActivePackageOutput struct {
	Packages []any `json:"packages"`
}

// SubmitQuestionnaireOutput output
type SubmitQuestionnaireOutput struct {
	ResultID uuid.UUID `json:"result_id"`
	Score    any       `json:"score"`
	TODO     []any     `json:"any"`
}

// GetChildStatOutput output
type GetChildStatOutput struct {
	TODO []any `json:"todo"`
}

// SearchQUestionnaireResultsOutput output
type SearchQUestionnaireResultsOutput struct {
	TODO []any `json:"todo"`
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
type UpdateChildernOutput struct{}

// GetMyChildernOutput output
type GetMyChildernOutput struct {
	ChildernData any `json:"childern_data"`
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
}

// SearchChildernOutput output
type SearchChildernOutput struct {
	Childern []any `json:"childern"`
}
