package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/delivery/rest"
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
