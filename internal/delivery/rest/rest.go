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
// @securityDefinitions.apikey	AdministratorLevelAuth
// @description				Bearer Token authentication for secure endpoints accessible only by administrator level user
// @in							header
// @name						Authorization
//
// @securityDefinitions.apikey	ParentLevelAuth
// @description				Bearer Token authentication for secure endpoints accessible only by registered user with auth token
// @in							header
// @name						Authorization

// @securityDefinitions.apikey	TherapistLevelAuth
// @description				Bearer Token authentication for secure endpoints accessible only by therapist level user
// @in							header
// @name						Authorization
type Service struct {
	v1                   *echo.Group
	authUsecase          usecase.AuthUsecaseIface
	packageUsecase       usecase.PackageUsecaseIface
	childUsecase         usecase.ChildUsecaseIface
	questionnaireUsecase usecase.QuestionnaireUsecaseIface
	usersUsecase         usecase.UsersUsecaseIface
}

// NewService init rest Service
func NewService(
	v1 *echo.Group, authUsecase usecase.AuthUsecaseIface,
	packageUsecase usecase.PackageUsecaseIface, childUsecase usecase.ChildUsecaseIface,
	questionnaireUsecase usecase.QuestionnaireUsecaseIface,
	usersUsecase usecase.UsersUsecaseIface,
) *Service {
	s := &Service{
		v1:                   v1,
		authUsecase:          authUsecase,
		packageUsecase:       packageUsecase,
		childUsecase:         childUsecase,
		questionnaireUsecase: questionnaireUsecase,
		usersUsecase:         usersUsecase,
	}

	s.initV1Routes()

	return s
}

func (s *Service) initV1Routes() {
	s.v1.POST("/auth/signup", s.HandleSignUp())
	s.v1.POST("/auth/signup/resend", s.HandleResendSignupVerification())
	s.v1.GET("/auth/verify", s.HandleVerifyAccount())
	s.v1.POST("/auth/login", s.HandleLogin())
	s.v1.PATCH("/auth/password", s.HandleInitResetPassword())
	s.v1.POST("/auth/password", s.HandleResetPassword())
	s.v1.GET("/auth/password", s.HandleRenderChangePasswordPage())
	s.v1.DELETE("/auth/accounts", s.HandleDeleteAccount(), s.AuthMiddleware(false))

	s.v1.POST("/atec/packages", s.HandleCreatePackage(), s.AuthMiddleware(false))
	s.v1.PUT("/atec/packages/:package_id", s.HandleUpdatePackage(), s.AuthMiddleware(false))
	s.v1.PATCH("/atec/packages/:package_id", s.HandleActivationPackage(), s.AuthMiddleware(false))
	s.v1.DELETE("/atec/packages/:package_id", s.HandleDeletePackage(), s.AuthMiddleware(false))
	s.v1.GET("/atec/packages/active", s.HandleSearchActivePackage())

	s.v1.POST("/childern", s.HandleRegisterChildern(), s.AuthMiddleware(false))
	s.v1.PUT("/childern/:child_id", s.HandleUpdateChildern(), s.AuthMiddleware(false))
	s.v1.GET("/childern", s.HandleGetMyChildern(), s.AuthMiddleware(false))
	s.v1.GET("/childern/search", s.HandleSearchChildern(), s.AuthMiddleware(false))
	s.v1.GET("/childern/:child_id/stats", s.HandleGetChildStats(), s.AuthMiddleware(false))

	s.v1.GET("/atec/questionnaires", s.HandleGetATECQuestionaire())
	s.v1.POST("/atec/questionnaires", s.HandleSubmitQuestionnaire(), s.AuthMiddleware(true))
	s.v1.GET(
		"/atec/questionnaires/results/:result_id",
		s.HandleDownloadQuestionnaireResult(), s.AuthMiddleware(true),
	)
	s.v1.GET("/atec/questionnaires/results", s.HandleSearchQUestionnaireResults(), s.AuthMiddleware(false))
	s.v1.GET("/atec/questionnaires/results/my", s.HandleGetMyQUestionnaireResults(), s.AuthMiddleware(false))

	s.v1.GET("/users/me", s.HandleGetMyProfile(), s.AuthMiddleware(false))

	s.v1.GET("/users/therapists", s.HandleGetTherapists(), s.AuthMiddleware(false))

	s.v1.GET("/swagger/*", echoSwagger.WrapHandler)
}

// StandardSuccessResponse base model for any successful API response
type StandardSuccessResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" example:"success"`
	Data       any    `json:"data"`
}

// StandardErrorResponse base model for any failed / error API response
type StandardErrorResponse struct {
	StatusCode   int    `json:"status_code" example:"400"`
	ErrorMessage string `json:"error_message" example:"Bad Request"`
	ErrorCode    string `json:"error_code" example:"missing required fields on input"`
}
