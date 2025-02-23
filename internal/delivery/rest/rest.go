// Package rest contain all functionality needed to use REST as interface to this API
package rest

import (
	"github.com/labstack/echo/v4"
	_ "github.com/luckyAkbar/atec/docs" // required by swaggo
	"github.com/luckyAkbar/atec/internal/usecase"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title						ATEC API Docs
// @version					1.0
// @description				this is the ATEC method implemented using API
// @securityDefinitions.apikey	AdminLevelAuth
// @description				Bearer Token authentication for secure endpoints accessible only by Admin level user
// @in							header
// @name						Authorization
//
// @securityDefinitions.apikey	UserLevelAuth
// @description				Bearer Token authentication for secure endpoints accessible only by registered user with auth token
// @in							header
// @name						Authorization
type service struct {
	v1                   *echo.Group
	authUsecase          usecase.AuthUsecaseIface
	packageUsecase       usecase.PackageUsecaseIface
	childUsecase         usecase.ChildUsecaseIface
	questionnaireUsecase usecase.QuestionnaireUsecaseIface
}

// NewService init rest service
func NewService(
	v1 *echo.Group, authUsecase usecase.AuthUsecaseIface,
	packageUsecase *usecase.PackageUsecase, childUsecase *usecase.ChildUsecase,
	questionnaireUsecase *usecase.QuestionnaireUsecase,
) {
	s := &service{
		v1:                   v1,
		authUsecase:          authUsecase,
		packageUsecase:       packageUsecase,
		childUsecase:         childUsecase,
		questionnaireUsecase: questionnaireUsecase,
	}

	s.initV1Routes()
}

func (s *service) initV1Routes() {
	s.v1.POST("/auth/signup", s.HandleSignUp())
	s.v1.POST("/auth/signup/resend", s.HandleResendSignupVerification())
	s.v1.GET("/auth/verify", s.HandleVerifyAccount())
	s.v1.POST("/auth/login", s.HandleLogin())
	s.v1.PATCH("/auth/password", s.HandleInitResetPassword())
	s.v1.POST("/auth/password", s.HandleResetPassword())
	s.v1.GET("/auth/password", s.HandleRenderChangePasswordPage())

	s.v1.POST("/atec/packages", s.HandleCreatePackage(), s.AuthMiddleware(false, true))
	s.v1.PUT("/atec/packages/:package_id", s.HandleUpdatePackage(), s.AuthMiddleware(false, true))
	s.v1.PATCH("/atec/packages/:package_id", s.HandleActivationPackage(), s.AuthMiddleware(false, true))
	s.v1.DELETE("/atec/packages/:package_id", s.HandleDeletePackage(), s.AuthMiddleware(false, true))
	s.v1.GET("/atec/packages/active", s.HandleSearchActivePackage())

	s.v1.POST("/childern", s.HandleRegisterChildern(), s.AuthMiddleware(true, false))
	s.v1.PUT("/childern/:child_id", s.HandleUpdateChildern(), s.AuthMiddleware(true, false))
	s.v1.GET("/childern", s.HandleGetMyChildern(), s.AuthMiddleware(true, false))
	s.v1.GET("/childern/search", s.HandleSearchChildern(), s.AuthMiddleware(false, true))
	s.v1.GET("/childern/:child_id/stats", s.HandleGetChildStats(), s.AuthMiddleware(true, false))

	s.v1.GET("/atec/questionnaires", s.HandleGetATECQuestionaire())
	s.v1.POST("/atec/questionnaires", s.HandleSubmitQuestionnaire(), s.allowUnauthorizedAccess(), s.AuthMiddleware(true, false))
	s.v1.GET(
		"/atec/questionnaires/results/:result_id",
		s.HandleDownloadQuestionnaireResult(),
		s.allowUnauthorizedAccess(), s.AuthMiddleware(true, false),
	)
	s.v1.GET("/atec/questionnaires/results", s.HandleSearchQUestionnaireResults(), s.AuthMiddleware(false, true))
	s.v1.GET("/atec/questionnaires/results/my", s.HandleGetMyQUestionnaireResults(), s.AuthMiddleware(true, false))

	s.v1.GET("/swagger/*", echoSwagger.WrapHandler)
}

// StandardSuccessResponse model info
//
//	@Description	base model for any successful API response
type StandardSuccessResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" example:"success"`
	Data       any    `json:"data"`
}

// StandardErrorResponse model info
//
//	@Description	base model for any failed / error API response
type StandardErrorResponse struct {
	StatusCode   int    `json:"status_code" example:"400"`
	ErrorMessage string `json:"error_message" example:"Bad Request"`
	ErrorCode    string `json:"error_code" example:"missing required fields on input"`
}
