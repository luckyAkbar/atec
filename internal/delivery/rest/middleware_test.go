package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	usecase_mock "github.com/luckyAkbar/atec/mocks/internal_/usecase"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/atec/internal/delivery/rest"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRESTMiddleware_AuthMiddleware(t *testing.T) {
	e := echo.New()
	group := e.Group("")

	mockAuthUsecase := usecase_mock.NewAuthUsecaseIface(t)

	service := rest.NewService(group, mockAuthUsecase, nil, nil, nil, nil)

	fn := func(allowUnauthorized bool) echo.HandlerFunc {
		return func(c echo.Context) error {
			if allowUnauthorized {
				return c.String(http.StatusOK, "ok")
			}

			user := model.GetUserFromCtx(c.Request().Context())
			if user == nil {
				t.Error("user is expected not to be nil if allowUnauthorized is false")
			}

			return c.String(http.StatusOK, "ok")
		}
	}

	testCases := []struct {
		name              string
		reqCtx            func() (*httptest.ResponseRecorder, echo.Context)
		expect            func(rec *httptest.ResponseRecorder)
		mockFn            func(ectx echo.Context)
		allowUnauthorized bool
	}{
		{
			name:              "disallow unauthorized, and missing token",
			allowUnauthorized: false,
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rec := httptest.NewRecorder()

				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name:              "disallow unauthorized, malformed token",
			allowUnauthorized: false,
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rec := httptest.NewRecorder()

				req.Header.Set("Authorization", " ")

				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name:              "allow unauthorized, and missing token just continue",
			allowUnauthorized: true,
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rec := httptest.NewRecorder()

				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name:              "disallow unauthorized, usecase returning error",
			allowUnauthorized: false,
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rec := httptest.NewRecorder()

				req.Header.Set("Authorization", "secretboss")

				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().AuthenticateAccessToken(ectx.Request().Context(), usecase.AuthenticateAccessTokenInput{
					Token: "secretboss",
				}).Return(nil, usecase.UsecaseError{
					ErrType: usecase.ErrUnauthorized,
				}).Times(1)
			},
		},
		{
			name:              "disallow unauthorized, ok",
			allowUnauthorized: false,
			reqCtx: func() (*httptest.ResponseRecorder, echo.Context) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rec := httptest.NewRecorder()

				req.Header.Set("Authorization", "secretboss")

				ectx := e.NewContext(req, rec)

				return rec, ectx
			},
			expect: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
			},
			mockFn: func(ectx echo.Context) {
				mockAuthUsecase.EXPECT().AuthenticateAccessToken(ectx.Request().Context(), usecase.AuthenticateAccessTokenInput{
					Token: "secretboss",
				}).Return(&usecase.AuthenticateAccessTokenOutput{
					UserID:   uuid.New(),
					UserRole: model.RolesAdministrator,
				}, nil).Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec, ectx := tc.reqCtx()

			if tc.mockFn != nil {
				tc.mockFn(ectx)
			}

			err := service.AuthMiddleware(tc.allowUnauthorized)(fn(tc.allowUnauthorized))(ectx)

			require.NoError(t, err)
			tc.expect(rec)
		})
	}
}
