package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/delivery/rest"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	usecase_mock "github.com/luckyAkbar/atec/mocks/internal_/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_HandleSignUp(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)
	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "validation error - email invalid",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{
					"email": "invalid-email",
					"password": "password12345",
					"username": "username"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleSignup(ectx.Request().Context(), usecase.SignupInput{
					Email:    "invalid-email",
					Password: "password12345",
					Username: "username",
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{
					"email": "valid@email.test",
					"password": "password12345",
					"username": "username"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleSignup(ectx.Request().Context(), usecase.SignupInput{
					Email:    "valid@email.test",
					Password: "password12345",
					Username: "username",
				}).Return(&usecase.SignupOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleSignUp()(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestAuthService_HandleResendSignupVerification(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)
	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/signup/resend", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "validation error - email invalid",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/signup/resend", strings.NewReader(`{
					"email": "invalid-email"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleResendSignupVerification(ectx.Request().Context(), usecase.ResendSignupVerificationInput{
					Email: "invalid-email",
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/signup/resend", strings.NewReader(`{
					"email": "valid@email.test",
					"password": "password12345",
					"username": "username"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleResendSignupVerification(ectx.Request().Context(), usecase.ResendSignupVerificationInput{
					Email: "valid@email.test",
				}).Return(&usecase.ResendSignupVerificationOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleResendSignupVerification()(ectx)
			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestAuthService_HandleVerifyAccount(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)
	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "validation error - token not found",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/auth/verify", nil)

				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleAccountVerification(ectx.Request().Context(), usecase.AccountVerificationInput{
					VerificationToken: "",
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/auth/verify", nil)

				q := req.URL.Query()
				q.Add("validation_token", "valid-token")
				req.URL.RawQuery = q.Encode()

				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)
				ectx.SetPath("/auth/verify")
				ectx.SetParamNames("token", "code")
				ectx.SetParamValues("valid-token", "valid-code")

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "text/html; charset=UTF-8", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleAccountVerification(ectx.Request().Context(), usecase.AccountVerificationInput{
					VerificationToken: "valid-token",
				}).Return(&usecase.AccountVerificationOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleVerifyAccount()(ectx)
			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestAuthService_HandleLogin(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)
	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "validation error - email invalid",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{
					"email": "invalid-email",
					"password": "password12345"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleLogin(ectx.Request().Context(), usecase.LoginInput{
					Email:    "invalid-email",
					Password: "password12345",
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{
					"email": "valid@email.test",
					"password": "password12345"
				}`))

				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)
				ectx.SetPath("/auth/login")

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleLogin(ectx.Request().Context(), usecase.LoginInput{
					Email:    "valid@email.test",
					Password: "password12345",
				}).Return(&usecase.LoginOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleLogin()(ectx)
			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestAuthService_HandleInitResetPassword(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)
	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPatch, "/auth/password", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "validation error - email invalid",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPatch, "/auth/password", strings.NewReader(`{
					"email": "invalid-email"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleInitesetPassword(ectx.Request().Context(), usecase.InitResetPasswordInput{
					Email: "invalid-email",
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPatch, "/auth/password", strings.NewReader(`{
					"email": "valid@email.test"
				}`))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)
				ectx.SetPath("/auth/password")
				ectx.SetParamNames("email", "password")
				ectx.SetParamValues("", "password12345")

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleInitesetPassword(ectx.Request().Context(), usecase.InitResetPasswordInput{
					Email: "valid@email.test",
				}).Return(&usecase.InitResetPasswordOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleInitResetPassword()(ectx)
			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestAuthService_HandleResetPassword(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)
	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/password", strings.NewReader(`{,}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "validation error - missing token",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/password", strings.NewReader(`{
					"token": "invalid-token",
					"new_password": "password12345"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleResetPassword(ectx.Request().Context(), usecase.ResetPasswordInput{
					NewPassword:        "password12345",
					ResetPasswordToken: "invalid-token",
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrBadRequest,
				}).Once()
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodPost, "/auth/password", strings.NewReader(`{
					"token": "valid-token",
					"new_password": "new-password"
				}`))

				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)
				ectx.SetPath("/auth/password")
				ectx.SetParamNames("token", "new_password")
				ectx.SetParamValues("valid-token", "new-password")

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleResetPassword(ectx.Request().Context(), usecase.ResetPasswordInput{
					NewPassword:        "new-password",
					ResetPasswordToken: "valid-token",
				}).Return(&usecase.ResetPasswordOutput{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleResetPassword()(ectx)
			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestAuthService_HandleRenderChangePasswordPage(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)
	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "validation error - missing change password token",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/auth/password", nil)

				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)
				ectx.SetPath("/auth/password")

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "ok",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodGet, "/auth/password", nil)

				q := req.URL.Query()
				q.Add(model.ChangePasswordTokenQuery, "valid-token")
				req.URL.RawQuery = q.Encode()

				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)
				ectx.SetPath("/auth/password")

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "text/html; charset=UTF-8", rec.Header().Get(echo.HeaderContentType))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleRenderChangePasswordPage()(ectx)
			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}

func TestAuthService_HandleDeleteAccount(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)
	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil)

	testCases := []struct {
		name   string
		reqCtx func() (*httptest.ResponseRecorder, echo.Context)
		expect func(rec *httptest.ResponseRecorder)
		mockFn func(ectx echo.Context)
	}{
		{
			name: "validation error - invalid input body",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodDelete, "/auth/account", strings.NewReader(`{
					"email": "",
					"password": "password12345",
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
		},
		{
			name: "usecase returning error",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodDelete, "/auth/account", strings.NewReader(`{
					"email": "valid@email.test",
					"password": "password12345"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleDeleteUserData(ectx.Request().Context(), usecase.DeleteUserDataInput{
					Email:    "valid@email.test",
					Password: "password12345",
				}).Return(usecase.UsecaseError{
					ErrType: usecase.ErrInternal,
				}).Once()
			},
		},
		{
			name: "success",
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				req := httptest.NewRequest(http.MethodDelete, "/auth/account", strings.NewReader(`{
					"email": "valid@email.test",
					"password": "password12345"
				}`))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, rec.Code)
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().HandleDeleteUserData(ectx.Request().Context(), usecase.DeleteUserDataInput{
					Email:    "valid@email.test",
					Password: "password12345",
				}).Return(nil).Once()
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.HandleDeleteAccount()(ectx)
			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}
