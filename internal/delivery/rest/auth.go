package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/usecase"
)

// @Summary		Create a new account to access user only resources within this system
// @Description	Use this API endpoint to create a new account
// @Tags			Authentication
// @Accept			json
// @Param			signup_input	body	SignupInput	true	"Login Credentials"
// @Produce		json
// @Success		200	{object}	StandardSuccessResponse{data=SignupOutput}	"Successful response"
// @Failure		400	{object}	StandardErrorResponse						"Bad request / validation error"
// @Router			/v1/auth/signup [post]
func (s *service) HandleSignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &SignupInput{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, StandardErrorResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "failed to parse input zzz",
				ErrorCode:    http.StatusText(http.StatusBadRequest),
			})
		}

		output, err := s.authUsecase.HandleSignup(c.Request().Context(), usecase.SignupInput{
			Email:    input.Email,
			Password: input.Password,
			Username: input.Username,
		})

		if err != nil {
			return usecaseErrorToRESTResponse(c, err)
		}

		return c.JSON(http.StatusOK, StandardSuccessResponse{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data: SignupOutput{
				Message: output.Message,
			},
		})
	}
}

// @Summary		Validate account after signup
// @Description	Confirmation method to ensure the email used when signup is active and owned by requester.
// @Description	If the confirmation token is valid, the account will be activated
// @Description	and will be able to be used on login. Otherwise, the opposite will happen.
// @Tags			Authentication
// @Accept			json
// @Param			validation_token	query	string	true	"validation token received from email after signup in form of a jwt token"
// @Produce		json
// @Success		200	{object}	StandardSuccessResponse{data=VerifyAccountOutput}	"Successful response"
// @Failure		400	{object}	StandardErrorResponse								"Bad request / validation error"
// @Failure		401	{object}	StandardErrorResponse								"invalid or expired verification token"
// @Router			/v1/auth/verify [get]
func (s *service) HandleVerifyAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Gain access to the system by authenticating using a registered account
// @Description	Use this endpoint to login with your username and password
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			login_input	body		LoginInput									true	"account detail such as email and password to log in"
// @Success		200			{object}	StandardSuccessResponse{data=LoginOutput}	"Successful response"
// @Failure		400			{object}	StandardErrorResponse						"Bad request"
// @Failure		401			{object}	StandardErrorResponse						"Authentication Failed"
// @Failure		500			{object}	StandardErrorResponse						"Internal Error"
// @Router			/v1/auth/login [post]
func (s *service) HandleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Initiate change password process for an active account
// @Description	If the user wants to change their password, use this API.
// @Description	when the request succeed, an email containing confirmation link will
// @Description	be sent to the account email.
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			init_change_password_input	body		InitResetPasswordInput									true	"the email of the account which password will be reset"
// @Success		200							{object}	StandardSuccessResponse{data=InitResetPasswordOutput}	"Successful response"
// @Failure		400							{object}	StandardErrorResponse									"Bad request"	"validation error"
// @Failure		500							{object}	StandardErrorResponse									"Internal Error"
// @Router			/v1/auth/password [patch]
func (s *service) HandleInitResetPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}

// @Summary		Change user's account password
// @Description	After using init reset password method, the resulting JWT token can be used here
// @Description	to authorize the password changes.
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			change_password_input	body		ResetPasswordInput									true	"the email of the account which password will be reset"
// @Success		200						{object}	StandardSuccessResponse{data=ResetPasswordOutput}	"Successful response"
// @Failure		400						{object}	StandardErrorResponse								"Bad request"	"validation error"
// @Failure		500						{object}	StandardErrorResponse								"Internal Error"
// @Router			/v1/auth/password [post]
func (s *service) HandleResetPassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotImplemented)
	}
}
