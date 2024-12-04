package rest

import "github.com/google/uuid"

type SearchActivePackageOutput struct {
	Packages []any `json:"packages"`
}

type SubmitQuestionnaireOutput struct {
	ResultID uuid.UUID `json:"result_id"`
	Score    any       `json:"score"`
	TODO     []any     `json:"any"`
}

type GetChildStatOutput struct {
	TODO []any `json:"todo"`
}

type SearchQUestionnaireResultsOutput struct {
	TODO []any `json:"todo"`
}

type GetMyQUestionnaireResultsOutput struct {
	TODO []any `json:"todo"`
}

type SignupOutput struct {
	Message string `json:"message" example:"confirmation link sent to your email"`
}

type VerifyAccountOutput struct {
	Message string `json:"message" example:"your account is now activated and can be used"`
}

type InitResetPasswordOutput struct {
}

type ResetPasswordOutput struct {
}

type LoginOutput struct {
	Token string `json:"token"`
}

type RegisterChildOutput struct {
}

type UpdateChildernOutput struct{}

type GetMyChildernOutput struct {
	ChildernData any `json:"childern_data"`
}

type CreatePackageOutput struct {
}

type UpdatePackageOutput struct {
}

type ActivationPackageOutput struct {
	Status bool `json:"status"`
}

type CreateATECTemplateOutput struct {
	ID uuid.UUID `json:"id"`
}

type UpdateATECTemplateOutput struct {
	ID uuid.UUID `json:"id"`
}

type ActivateTemplateOutput struct {
	Status bool `json:"status"`
}

type GetATECQuestionnaireOutput struct {
}
